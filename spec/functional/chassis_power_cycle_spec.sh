#!/usr/bin/env sh
# MIT License
#
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP
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

Describe "gru --config ${GRU_CONF} chassis power"
BeforeAll use_valid_config
BeforeAll use_valid_bios_attributes_file

# test all "cycle" variant against all vendors/models as their outputs vary
# (see chassis power --help and testdata/fixtures/rie)
Parameters:matrix
  "cycle" "reset"
  "127.0.0.1:5000" "127.0.0.1:5001" "127.0.0.1:5002" "127.0.0.1:5003" "127.0.0.1:5004"
End

# check that the requested power state is on
It "$1 $2"
  When call ./gru --config "${GRU_CONF}" chassis power "$1" "$2"
  The status should equal 0
  The line 1 of stdout should include "$2:" # the host should be in the stdout
  # the power state may vary, but check for the word 'PreviousPowerState'
  The line 2 of stdout should include 'PreviousPowerState' 
  # Powering on should also show the requested power state
  The line 3 of stdout should include 'RequestedPowerState'
  The lines of stderr should equal 1
End

# validate yaml and json outputs work
It "$1 $2 --yaml"
  When call ./gru --config "${GRU_CONF}" chassis power "$1" "$2" "--yaml"
  The status should equal 0
  The stderr should be present
  The stdout should "be_yaml"
End
# FIXME: newlines in JSON with invalid reset types: 
#        "message": "{\n    \"Status\": 400,\n    \"Message\": \"Invalid ResetType\"\n}",
#        jq: parse error: Invalid string: control characters from U+0000 through U+001F must be escaped at line 10, column 2
# It "$1 $2 --json"
#   When call ./gru  --config "${GRU_CONF}" chassis power "$1" "$2" "--json"
#   The status should equal 0
#   The stderr should be present
#   The stdout should "be_json"
# End

# validate piping to STDIN works
Data:expand
 #| $2
End
It "$1 $2 (host passed via STDIN)"
  When call ./gru  --config "${GRU_CONF}" chassis power "$1"
  The status should equal 0
  The stderr should be present
  The stdout should be present
End

End
