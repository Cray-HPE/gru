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

# test various combinations of --attributes, which may or may not exist from a pariticlar vendor/model
# test comma-separated list of attributes as well
# test all vendors/models (see testdata/fixtures/rie) as their outputs vary
Parameters
  127.0.0.1:5000   
  127.0.0.1:5001              
  # 127.0.0.1:5002   
  127.0.0.1:5003      
  # 127.0.0.1:5004   
End

# getting a single key should return only that key (1 line of stdout)
It "$1 --attributes SingleKey"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --attributes "SingleKey"
  The status should equal 0
  The line 1 of stdout should include "$1:"
  The line 2 of stdout should include 'Attributes'
  The line 3 of stdout should include 'SingleKey'
  The lines of stdout should equal 3
  The lines of stderr should equal 1
End

# getting multiple keys should return only those keys
It "$1 --attributes Key1,Key2"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --attributes  Key1,Key2
  The status should equal 0
  The line 1 of stdout should include "$1:"
  The line 2 of stdout should include 'Attributes'
  The line 3 of stdout should include 'Key1'
  The line 4 of stdout should include 'Key2'
  The lines of stdout should equal 4
  The lines of stderr should equal 1
End

# validate yaml and json outputs work
It "$1 --attributes SingleKey --yaml"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --attributes "SingleKey" "--yaml"
  The status should equal 0
  The stderr should be present
  The stdout should "be_yaml"
End
It "$1 --attributes SingleKey --json"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --attributes "SingleKey" "--json"
  The status should equal 0
  The stderr should be present
  The stdout should "be_json"
End

# validate piping to STDIN works
Data:expand
 #| $1
End
It "--attributes SingleKey (host passed via STDIN)"
  When call ./gru --config "${GRU_CONF}" bios get --attributes SingleKey
  The status should equal 0
  The stderr should be present
  The stdout should be present
End

End

# Vendor-specific tests, mostly relating to those vendors/models needing a decoder for their BIOS attribute names
Describe "gru --config ${GRU_CONF} bios get"
Parameters
  # $1           $2         $3                 $4           $5                  $6           $7               $8            $9         $10           $11           $12
  127.0.0.1:5001 "PCIS007"  "SR-IOV Support"   "Rome0039"   "Local APIC Mode"   "Rome0059"   "SMT Control"    "Rome0162"    "IOMMU"    "Rome0565"    "SVM Mode"    "7"
  # add other vendors if they have a decoder
End

# --virtualization shortcut should return only virtualization attributes and show the code names and friendly names
It "$1 --virtualization"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --virtualization
  The status should equal 0
  The stdout should include "${2} (${3})"
  The stdout should include "${4} (${5})"
  The stdout should include "${6} (${7})"
  The stdout should include "${8} (${9})"
  The stdout should include "${10} (${11})"
  The lines of stdout should equal "${12}"
  The lines of stderr should equal 1
End

# --virtualization shortcut should return only virtualization attributes and only show the code names in the json output
It "$1 --virtualization --json"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --virtualization --json
  The status should equal 0
  The stdout should include "${2}"
  The stdout should include "${4}"
  The stdout should include "${6}"
  The stdout should include "${8}"
  The stdout should include "${10}"
  The lines of stderr should equal 1
  The stdout should be_json
End

# --virtualization shortcut should return only virtualization attributes and only show the code names in the yaml output
It "$1 --virtualization --yaml"
  When call ./gru --config "${GRU_CONF}" bios get "$1" --virtualization --yaml
  The status should equal 0
  The stdout should include "${2}"
  The stdout should include "${4}"
  The stdout should include "${6}"
  The stdout should include "${8}"
  The stdout should include "${10}"
  The lines of stderr should equal 1
  The stdout should be_yaml
End

End
