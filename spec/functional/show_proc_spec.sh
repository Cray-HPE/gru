#!/usr/bin/env sh
# MIT License
#
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP
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


Describe "gru show proc --config ${GRU_CONF}"
BeforeAll use_valid_config
BeforeAll use_valid_bios_attributes_file

# test all vendors/models (see testdata/fixtures/rie) as their outputs vary
Parameters
  127.0.0.1:5000    "CPU 1"     "CPU 2"   
  127.0.0.1:5001    "P1"        "P1"      
  127.0.0.1:5002    "Proc 1"    ""       
  127.0.0.1:5003    "CPU 0"     "CPU 1"   
  127.0.0.1:5004    "Proc 1"    "Proc 2"  
End

# the proc names vary by vendor/model, so check them according to the parameters
It "$1"
  When call ./gru --config "${GRU_CONF}" show proc "$1"
  The status should equal 0
  The stdout should include "$2"
  The stdout should include "$3"
  The lines of stderr should equal 1
End

# validate yaml and json outputs work
It "$1 --yaml"
  When call ./gru --config "${GRU_CONF}" show proc "$1" "--yaml"
  The status should equal 0
  The stderr should be present
  The stdout should "be_yaml"
End
It "$1 --json"
  When call ./gru --config "${GRU_CONF}" show proc "$1" "--json"
  The status should equal 0
  The stderr should be present
  The stdout should "be_json"
End

# also check that stdin works for this command, just checking that output exists
Data:expand
 #| $1
End
It "$1 (host passed via STDIN)"
  When call ./gru --config "${GRU_CONF}" show proc
  The status should equal 0
  The stderr should be present
  The stdout should be present
End

End
