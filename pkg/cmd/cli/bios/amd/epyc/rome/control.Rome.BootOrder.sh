#set -x
log=admin
pass=password
bmc=$1

function help() {
  echo "Usage: Manages the BootOrder at /redfish/v1/Systems/Self"
  echo "# sh $0 < bmc > < show_current | update_random | update_manual boot1,boot2,boot4,..,boot0 >"
  echo
  echo "for examples"
  echo "# sh $0 r272-z30-mgmt show_current                                                 # Display current BootOrder and their descriptione"
#  echo "# sh $0 r272-z30-mgmt show_pending                                                 # Display current BootOrder and their descriptione"
  echo "# sh $0 r272-z30-mgmt update_random                                                # Change current BootOrder to Random BootOrder"
  echo "# sh $0 r272-z30-mgmt update_manual Boot0003,Boot0002,X,X,Boot0000                 # Update BootOrder in order of entered boot devices."
  exit 1
}

function show_current() {
  printf "[ BMC: ${bmc} ]\n"
  printf "ex: curl -ksu $log:$pass -H \"content-type: application/json\" https://${bmc}/redfish/v1/Systems/Self | jq '.Boot .BootOrder | .[]' | sed 's/\"//g'\n"
  printf "ex: curl -ksu $log:$pass -H \"content-type: application/json\" https://${bmc}/redfish/v1/Systems/Self/BootOptions/000C | jq '.Description'\n"
  bootorderlist=$(curl -ksu $log:$pass -H "content-type: application/json" https://${bmc}/redfish/v1/Systems/Self | jq '.Boot .BootOrder | .[]' | sed 's/"//g')
  for i in $(printf "$bootorderlist\n" | tr '\n' ' ') ; do
    bootdevice="$(printf "$i\n" | sed 's/Boot//g')"
    devicedesc=$(curl -ksu $log:$pass -H "content-type: application/json" https://${bmc}/redfish/v1/Systems/Self/BootOptions/${bootdevice} | jq '.Description')
    printf "%-10s %-30s\n" "Boot${bootdevice} -" "$devicedesc"
  done
}

function update_random() {
  printf "[ BMC: ${bmc} ]\n"
  randomorder=$(curl -ksu $log:$pass -H "content-type: application/json" https://${bmc}/redfish/v1/Systems/Self | jq '.Boot .BootOrder | .[]' | sed 's/"//g' | shuf)
  random=$(printf "$randomorder\n" | sed 's/^/"/g' | sed 's/$/"/g' | tr '\n' ',' | sed 's/,$//g') 						# change randomorder to "Boot0004","Boot0003","Boot0005" format
  printf "ex: curl -sku $log:$pass -H \"content-type: application/json\" -H \"If-Match: *\" -X PATCH https://${bmc}/redfish/v1/Systems/Self -d '{\"Boot\":{\"BootOrder\":[$random]}}' | jq\n"
  curl -sku $log:$pass -H "content-type: application/json" -H "If-Match: *" -X PATCH https://${bmc}/redfish/v1/Systems/Self -d '{"Boot":{"BootOrder":['$random']}}' | jq
  printf "Wait for 2 seconds, so Redfish could have new BootOrder.\n" ; sleep 2
  printf "[ New BootOrder ]\n"
  show_current
}

function update_manual() {
  newbootorder="$(printf "$1\n" | sed 's/^/"/g' | sed 's/$/"/g' | sed 's/,/","/g')"
  printf "ex: curl -sku $log:$pass -H \"content-type: application/json\" -H \"If-Match: *\" -X PATCH https://${bmc}/redfish/v1/Systems/Self -d '{\"Boot\":{\"BootOrder\":[$newbootorder]}}' | jq\n"
  curl -sku $log:$pass -H "content-type: application/json" -H "If-Match: *" -X PATCH https://${bmc}/redfish/v1/Systems/Self -d '{"Boot":{"BootOrder":['$newbootorder']}}' | jq
  printf "Wait for 5 seconds, so Redfish could have new BootOrder.\n" ; sleep 5
  printf "[ New BootOrder ]\n"
  bootorderlistpending=$(curl -ksu $log:$pass -H "content-type: application/json" https://${bmc}/redfish/v1/Systems/Self | jq '.Boot .BootOrder | .[]' | sed 's/"//g')
  for i in $(printf "$bootorderlistpending\n" | tr '\n' ' ') ; do
    bootdevice="$(printf "$i\n" | sed 's/Boot//g')"
    devicedesc=$(curl -ksu $log:$pass -H "content-type: application/json" https://${bmc}/redfish/v1/Systems/Self/BootOptions/${bootdevice} | jq '.Description')
    printf "%-10s %-30s\n" "Boot${bootdevice} -" "$devicedesc"
  done
}

[ "$2" == "" ] || [ "$2" != "show_current" ] && [ "$2" != "update_random" ] && [ "$2" != "update_manual" ] && help
[ "$2" == "show_current" ] && show_current
[ "$2" == "update_random" ] && update_random
[ "$2" == "update_manual" ] && update_manual $3
