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


Describe 'gru chassis power cycle'

# Running against an active host with good credentials should succeed and show output
It "--config ${GRU_CONF} 127.0.0.1:5000"
  BeforeCall use_valid_config
  When call ./gru chassis power cycle --config "${GRU_CONF}" 127.0.0.1:5000
  The status should equal 0
  The line 1 of stderr should include 'Asynchronously updating'
  The line 1 of stdout should equal '127.0.0.1:5000:'
  The line 2 of stdout should include 'PreviousPowerState'
  The line 3 of stdout should include 'RequestedPowerState'
  The line 3 of stdout should include 'GracefulRestart'
  The lines of stderr should equal 1
End
  
# validate piping to STDIN works
It "--config ${GRU_CONF} (via STDIN)"
  BeforeCall use_valid_config
  
  Data "127.0.0.1:5000" # STDIN

  When call ./gru chassis power cycle --config "${GRU_CONF}"
  The status should equal 0
  The line 1 of stderr should include 'Asynchronously updating'
  The line 1 of stdout should equal '127.0.0.1:5000:'
  The line 2 of stdout should include 'PowerState'
  The line 3 of stdout should include 'GracefulRestart'
  The line 4 of stdout should include "reset type 'GracefulRestart' is not supported by this service"
  The lines of stderr should equal 1
End

End
