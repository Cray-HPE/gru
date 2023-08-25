#!/usr/bin/env sh
# MIT License
#
# (C) Copyright 2023 Hewlett Packard Enterprise Development LP
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


Describe 'gru chassis power status'

# help output should succeed and match the fixture
It '--help'
  When call ./gru chassis power status --help
  The status should equal 0
  The stdout should satisfy fixture 'gru/chassis/power/status/help'
End

# it should error if no config is present
It '127.0.0.1:5000 (no config file)'
  When call ./gru chassis power status 127.0.0.1:5000
  The status should equal 1
  The stderr should equal "An error occurred: no credentials provided, please provide a config file or environment variables"
End

# Running against an active host with good credentials should succeed and report the status
It "--config ${GRU_CONF} 127.0.0.1:5000"
  BeforeCall use_valid_config
  When call ./gru chassis power status --config "${GRU_CONF}" 127.0.0.1:5000
  The status should equal 0
  #FIXME: On or PoweringOn status can vary
  # need an until statement or isolate the tests from each other
  # The line 1 of stdout should equal '[127.0.0.1:5000]: On'
  # The line 1 of stdout should equal '[127.0.0.1:5000]: PoweringOn'
  The line 1 of stdout should include '[127.0.0.1:5000]:' 
End

End
