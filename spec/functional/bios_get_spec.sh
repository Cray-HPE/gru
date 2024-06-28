#!/usr/bin/env sh
# MIT License
#
# (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP
#
# Permissioff is hereby granted, free of charge, to any persoff obtaining a
# copy of this software and associated documentatioff files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permissioff notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTIoff OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTIoff WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.


Describe "gru --config ${GRU_CONF} bios get"
BeforeAll use_valid_config
BeforeAll use_valid_bios_attributes_file

# test all vendors running in different containers with different ports (see testdata/fixtures/rie)
Parameters
  # Hostname:Port    
  # $1               $2      $3                                         $4     $5                      $6
  127.0.0.1:5000     "866"   '"Attributes" does not exist or is null'   "2"    "ProcessorVmxEnable"    "ProcessorX2apic"
  127.0.0.1:5001     "552"   "The resource at the URI"                  "19"   ""                      ""
  # 127.0.0.1:5002  
  127.0.0.1:5003     "20"    "The resource at the URI"                  "19"   ""                      ""
  # 127.0.0.1:5004
End

# getting all bios keys should return lots of output and will vary by vendor/model
It "$1"
  When call ./gru --config "${GRU_CONF}" bios get "$1"
  The status should equal 0
  The lines of stdout should equal "${2}" 
  The lines of stderr should equal 1
End

# getting pending changes should return an error if the Bios/Settings.Attributes does not exist
It "$1 --pending"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --pending 
  The status should equal 0
  The stdout should include "${3}"
  The lines of stdout should equal "${4}"
  The lines of stderr should equal 1
End

# validate yaml and json outputs work
It "$1 --pending --yaml"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --pending --yaml
  The status should equal 0
  The stderr should be present
  The stdout should "be_yaml"
End
It "$1 --pending --json"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --pending --json
  The status should equal 0
  The stderr should be present
  The stdout should "be_json"
End

# getting keys from a file should return those keys
It "$1 --from-file ${GRU_BIOS_KV}"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --from-file "${GRU_BIOS_KV}"
  The status should equal 0
  The stdout should include 'BootTimeout'
  The stdout should include 'SRIOVEnable'
  The lines of stdout should equal 4
  The lines of stderr should equal 1
End

# validate yaml and json outputs work
It "$1 --from-file ${GRU_BIOS_KV} --yaml"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --from-file "${GRU_BIOS_KV}" --yaml
  The status should equal 0
  The stderr should be present
  The stdout should "be_yaml"
End
It "$1 --from-file ${GRU_BIOS_KV} --json"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --from-file "${GRU_BIOS_KV}" --json
  The status should equal 0
  The stderr should be present
  The stdout should "be_json"
End


# re-enable after https://github.com/Cray-HPE/csm-redfish-interface-emulator/pull/10 merges
# # passing a shortcut should return a limited set of pre-defined keys
# It "$1 --virtualization"
#   When call ./gru bios get --config "${GRU_CONF}" "$1" --virtualization
#   The status should equal 0
#   The stdout should include "${5}"
#   The stdout should include "${6}"
#   The lines of stderr should equal 1
# End

# # validate yaml and json outputs work
# It "$1 --virtualization --yaml"
#   When call ./gru --config "${GRU_CONF}" bios get "$1" --virtualization --yaml
#   The status should equal 0
#   The stderr should be present
#   The stdout should "be_yaml"
# End
# It "$1 --virtualization --json"
#   When call ./gru --config "${GRU_CONF}" bios get "$1" --virtualization --json
#   The status should equal 0
#   The stderr should be present
#   The stdout should "be_json"
# End


# also check that stdin works for this command, just checking that output exists
Data:expand
 #| $1
End
It "$1 --pending (host passed via STDIN)"
  When call ./gru --config "${GRU_CONF}" bios get
  The status should equal 0
  The stderr should be present
  The stdout should be present
End

End
