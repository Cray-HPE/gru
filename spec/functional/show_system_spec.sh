#!/usr/bin/env sh
# MIT License
#
# (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.


Describe "gru --config ${GRU_CONF} show system"
BeforeAll use_valid_config
BeforeAll use_valid_bios_attributes_file

# test all vendors/models (see testdata/fixtures/rie) as their outputs vary
Parameters
  # $1              $2
  127.0.0.1:5000    "S2600BPB"                   
  127.0.0.1:5001    "H262-Z63-00"                
  127.0.0.1:5002    "ProLiant DL325 Gen10 Plus" 
  127.0.0.1:5003    "HPE CRAY EX425 (MILAN)"     
  127.0.0.1:5004    "ProLiant XL675d Gen10 Plus" 
End

# check the model number as it will vary by vendor/model
It "$1"
  When call ./gru --config "${GRU_CONF}" show system "$1"
  The status should equal 0
  The stdout should include "$2"
  The lines of stderr should equal 1
End

# validate yaml and json outputs work
It "$1 --yaml"
  When call ./gru --config "${GRU_CONF}" show system "$1" "--yaml"
  The status should equal 0
  The stderr should be present
  The stdout should "be_yaml"
End
It "$1 --json"
  When call ./gru --config "${GRU_CONF}" show system "$1" "--json"
  The status should equal 0
  The stderr should be present
  The stdout should "be_json"
End

# also check that stdin works for this command, just checking that output exists
Data:expand
 #| $1
End
It "$1 (host passed via STDIN)"
  When call ./gru --config "${GRU_CONF}" show system
  The status should equal 0
  The stderr should be present
  The stdout should be present
End

End
