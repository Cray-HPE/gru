#       Revisions:
#       DATE            MODIFIED BY     SUMMARY
#       6/01/2023       Jae Byun        Changed printf to echo to handel '-', '%' text.
# 	7/05/2023 	Jae Byun	update_biosparameter routine to support the ${para} without double quotes.
#       7/10/2023       Jae Byun        Added [ "$4" != "true" ] && [ "$4" != "false" ]
#
#set -x
log=admin
pass=password
bmc=$1
option=$2

function help() {
  echo "Usage: Manages the BIOS parameters"
  echo "# sh $0 < bmc > < renew_json | read [AttributeName] | update AttributeName new_parameter | read_future | load_BIOS_Default >"
  echo
  echo "for examples"
  echo "# sh $0 vp2-mgmt renew_json                                   # Execute renew_json first for next options - Read /redfish/v1/Registries/BiosAttributeRegistry.json and generate json files for possible parameters"
  echo "# sh $0 vp2-mgmt read                                         # Read current BIOS parameters all"
  echo "# sh $0 vp2-mgmt read EagleStream0001                         # read specific AttributeName alone"
  echo "# sh $0 vp2-mgmt update EagleStream0001 \"Single LP\"           # Update new_parameter to specific AttributeName"
  echo "# sh $0 vp2-mgmt read_future                                  # Read Future BIOS parameter, Redfish endpoint is /redfish/v1/Systems/Self/Bios/SD"
  echo "# sh $0 vp2-mgmt Admin_password current_pass new_pass         # Set Administrator password (SETUP001)"
  echo "# sh $0 vp2-mgmt User_password current_pass new_pass          # Set User password (SETUP002)"
  echo "# sh $0 vp2-mgmt load_BIOS_Default                            # Load Default BIOS Parameter, reboot is required after this command"
  echo 
  echo "Note: Some BIOS parameters have a dependency on others, and this script is not prepared."
  exit 1
}

function renew_json() {
  printf "[ BMC: ${bmc} ]\n"
  rm -rf *.json ; n=0
  data1=$(curl -ksu $log:$pass -H 'content-type:application/json' https://${bmc}/redfish/v1/Registries/BiosAttributeRegistry.json | jq '.RegistryEntries .Attributes[]')
  data2=$(curl -ksu $log:$pass -H 'content-type:application/json' https://${bmc}/redfish/v1/Registries/BiosAttributeRegistry.json | jq '.RegistryEntries .Dependencies[]') 
  echo "$data1" | while IFS= read -r line ; do
#    if [ $(printf "$line\n" | grep -c "^{") -eq 1 ] ; then
    if [ $(echo "$line" | grep -c "^{") -eq 1 ] ; then
      ((n=n+1))
#      printf "$line\n" > ${n}.json
      echo "$line" > ${n}.json
      echo "n:$n"
    else
#      printf "$line\n" >> ${n}.json
      echo "$line" >> ${n}.json
    fi
#  done <<< $(curl -ksu $log:$pass -H 'content-type:application/json' https://${bmc}/redfish/v1/Registries/BiosAttributeRegistry.json | jq '.RegistryEntries .Attributes[]')
  done

  nn=$(echo "$data1" | grep -c AttributeName)
  n=1
  while [ $n -le $nn ] ; do 
    printf "[ file: $n.json ] - "
    newname=$(cat $n.json | grep AttributeName | awk '{print $2}' | sed 's/"//g' | sed 's/,//g')
    printf "changing $n.json file to $newname.json\n"
    mv $n.json $newname.json
    ((n=n+1))
  done
# ---------------------------- generate DependencyFor.json files ---------------------------------------------
  n=0
  echo "$data2" | while IFS= read -r line ; do
    if [ $(echo "$line" | grep -c "^{") -eq 1 ] ; then
      ((n=n+1))
      echo "n:$n"
      echo "$line" > ${n}.json
    else
      echo "$line" >> ${n}.json
    fi
#  done <<< $(curl -ksu $log:$pass -H 'content-type:application/json' https://${bmc}/redfish/v1/Registries/BiosAttributeRegistry.json | jq '.RegistryEntries .Dependencies[]')
  done
  nn=$(echo "$data2" | grep -c "Dependency\"")
  n=1
  while [ $n -le $nn ] ; do
    printf "[ file: $n.json ] - "
    newname1=$(cat $n.json | grep MapFromAttribute | awk '{print $2}' | sed 's/"//g' | sed 's/,//g' | uniq)
    newname2=$(cat $n.json | grep MapToAttribute   | awk '{print $2}' | sed 's/"//g' | sed 's/,//g')
    printf "changing $n.json file to Dependency.$newname1.$newname2.json\n"
    mv $n.json Dependency.$newname1.$newname2.json
    ((n=n+1))
  done

}

function read_biosparameter() {
  printf "[[ BMC: ${bmc} ]]\n"
  printf "%-18s %-49s %-37s %-37s %-20s\n" "[ AttributeName ]" "[ DisplayName ]" "[ CurrentValue ]" "[ DefaultValue ]" "[ Available \"Values\" ]"
  cat .current.txt | sed $'s/[^[:print:]\t]//g' | grep -v MAPIDS | while IFS= read -r line ; do
    d1=$(echo "$line" | awk '{print $1}' | sed 's/"//g')									# AttributeName	
    data=$(cat $d1.json)									
    d2=$(jq '.DisplayName'                     		<<< $data 2>/dev/null) 							# DisplayName
    d3=$(echo "$line" | awk '{print $2,$3,$4,$5,$6,$7}' | sed 's/[[:blank:]]*$//g')			 			# CurrentValue
    d4=$(jq '.DefaultValue'                    		<<< $data 2>/dev/null) 							# DefaultValue
    [ "$d3" != "$d4" ] && { echo ; printf "%-18s %-49s\n" 	" " " >>> [ CurrentValue:$d3 and DefaultValue:$d4 are not the same ] <<< " ; }
    temp1=$(jq '.Value[] | .ValueDisplayName' 		<<< $data 2>/dev/null)
    temp2=$(jq '.Value[] | .ValueName'        		<<< $data 2>/dev/null)
    [ "$temp1" != "$temp2" ] && { echo ; printf "%-18s %-49s\n" " " "  >>> [ ValueDisplayName and ValueName are not the same, please check $d1.json file ] <<<\n" ; }
    d5=$(jq '.Value[] | .ValueName'            		<<< $data 2>/dev/null | tr '\n' ' ')					# Available "Value(s)"	
    d5n=$(jq '.Value[] | .ValueName'            	<<< $data 2>/dev/null | wc -l)				
    if [ $d5n -lt 11 ] ; then
      printf "%-18s %-49s %-37s %-37s %-20s\n" "$d1" "$d2" "$d3" "$d4" "$d5"
    else
      printf "%-18s %-49s %-37s %-37s\n" "$d1" "$d2" "$d3" "$d4"
      printf "%-18s %-49s\n" " " "[ Available values ] - $d5"
    fi
#   done <<< $(cat .current.txt | sed $'s/[^[:print:]\t]//g' | grep -v MAPIDS)      # .curent.txt has something sepcial texts.
    done
}

function update_biosparameter() {
  double_quotes=$(cat ${attrib}.json | grep DefaultValue | perl -ne 'printf "%d\n", tr /"/" /')            # count the double quotes (") from the DefaultValue section
  [ $double_quotes -ne 2 ] && { \
    printf "ex: curl -ksu $log:$pass -H 'content-type: application/json' -H 'If-Match: *' -X PATCH https://${bmc}/redfish/v1/Systems/Self/Bios/SD -d '{\"Attributes\":{\"${attrib}\":\"${para}\"}}' | jq\n"
    curl -ksu $log:$pass -H 'content-type: application/json' -H 'If-Match: *' -X PATCH https://${bmc}/redfish/v1/Systems/Self/Bios/SD -d '{"Attributes":{"'"${attrib}"'":"'"${para}"'"}}' | jq ; }
  [ $double_quotes -eq 2 ] && { \
    temp=$(printf "$para\n" | sed 's/"//g') ; para=$temp
    printf "ex: curl -ksu $log:$pass -H 'content-type: application/json' -H 'If-Match: *' -X PATCH https://${bmc}/redfish/v1/Systems/Self/Bios/SD -d '{\"Attributes\":{\"${attrib}\":${para}}}' | jq\n"
    curl -ksu $log:$pass -H 'content-type: application/json' -H 'If-Match: *' -X PATCH https://${bmc}/redfish/v1/Systems/Self/Bios/SD -d '{"Attributes":{"'"${attrib}"'":'${para}'}}' | jq ; }
}

function read_future() {
  printf "ex: curl -ksu $log:$pass -H 'content-type: application/json' https://${bmc}/redfish/v1/Systems/Self/Bios/SD | jq '.Attributes'\n"  
  curl -ksu $log:$pass -H 'content-type: application/json' https://${bmc}/redfish/v1/Systems/Self/Bios/SD | jq '.Attributes'
}

function load_BIOS_Default() {
#  data=$(curl -ksu $log:$pass https://${bmc}/redfish/v1/UpdateService/FirmwareInventory/BMCImage1 | jq '.Version' | sed 's/"//g' | sed 's/\.//g') 			# Collect BMC Version, and change it as example 1.11.00 to 11100.
#  [ $data -lt 11100 ] && { \
    printf "ex: curl -ksu $log:$pass -H 'content-type: application/json' -X POST https://${bmc}/redfish/v1/Systems/Self/Bios/Actions/Bios.ResetBios -d '{\"ResetType\":\"Reset\"}' | jq\n"
    curl -ksu $log:$pass -H 'content-type: application/json' -X POST https://${bmc}/redfish/v1/Systems/Self/Bios/Actions/Bios.ResetBios -d '{"ResetType":"Reset"}' | jq
#  [ $data -le 11100 ] && { \
#    printf "ex: curl -ksu $log:$pass -H 'content-type: application/json' -X POST https://${bmc}/redfish/v1/Systems/Self/Bios/Actions/Bios.ResetBios -d '{}' | jq\n"
#    curl -ksu $log:$pass -H 'content-type: application/json' -X POST https://${bmc}/redfish/v1/Systems/Self/Bios/Actions/Bios.ResetBios -d '{}' | jq ; }
}

function set_password() {
  printf "ex: curl -ksu $log:$pass -H 'content-type: application/json' -X POST https://${bmc}/redfish/v1/Systems/Self/Bios/Actions/Bios.ChangePassword -d '{\"PasswordName\":\"${1}\",\"OldPassword\":\"${2}\",\"NewPassword\":\"${3}\"}'\n"
  curl -ksu $log:$pass -H 'content-type: application/json' -X POST https://${bmc}/redfish/v1/Systems/Self/Bios/Actions/Bios.ChangePassword -d '{"PasswordName":"'"${1}"'","OldPassword":"'"${2}"'","NewPassword":"'"${3}"'"}'
}

[ "$option" == "" ] || [ "$option" != "renew_json" ] && [ "$option" != "read" ] && [ "$option" != "update" ] && [ "$option" != "load_BIOS_Default" ] && [ "$option" != "read_future" ] && \
[ "$option" != "Admin_password" ] && [ "$option" != "User_password" ] && { help ; }
[ "$2" == "renew_json" ] && { renew_json ; }
[ "$2" == "read" ] && [ "$3" == "" ] && { curl -ksu $log:$pass -H 'content-type:application/json' https://${bmc}/redfish/v1/Systems/Self/Bios | jq '.Attributes' | grep [A-Za-z] | sed 's/://g' | sed 's/,//g' > .current.txt ; read_biosparameter ; }
[ "$2" == "read" ] && [ "$3" != "" ] && { curl -ksu $log:$pass -H 'content-type:application/json' https://${bmc}/redfish/v1/Systems/Self/Bios | jq '.Attributes' | grep [A-Za-z] | sed 's/://g' | sed 's/,//g' \
      | grep "$3" > .current.txt ; read_biosparameter ; }
[ "$2" == "update" ] && [ "$3" != "" ] && [ "$4" != "" ] && { [ ! -f ${3}.json ] && printf "AttributeName:${3} does not exist\n" && exit 1 ; \
  [ $(ls -al Dependency*.json | grep -c ${3}) -ne 0 ] && printf "${3} has some dependency with others, check these files\n" && ls -al Dependency*${3}*.json && exit 1 ; \
  [ $(cat ${3}.json | grep -c "${4}") -eq 0 ] && [ "$4" != "true" ] && [ "$4" != "false" ] && printf "new_parameter:${4} does not exist\n" && exit 1 ; attrib="${3}" ; para="${4}" ; update_biosparameter ; }
[ "$2" == "read_future" ] && { read_future ; }
[ "$2" == "load_BIOS_Default" ] && { load_BIOS_Default ; }
[ "$2" == "Admin_password" ] || [ "$2" == "User_password" ] && { [ "$2" == "Admin_password" ] && opt=SETUP001 ; [ "$2" == "User_password" ] && opt=SETUP002 ; current="$3" ; new="$4" ; set_password "$opt" "$current" "$new" ; }

# Example command to load the BIOS configuration file.
# curl -k -u admin:superuser -X POST https://10.14.53.239/redfish/v1/Systems/Self/Bios/SD -H 'Content-Type: application/json' -H 'Expect:' -d @/FW/BIOScfg-fw1_11-v1_0.json
