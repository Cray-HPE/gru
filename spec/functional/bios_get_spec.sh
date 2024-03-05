#!/usr/bin/env sh
# MIT License
#
# (C) Copyright 2023 Hewlett Packard Enterprise Development LP
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


Describe 'gru bios get'

BeforeAll use_valid_config
BeforeAll use_valid_bios_attributes_file

# getting all bios keys should return lots of output
It "--config ${GRU_CONF} 127.0.0.1:5000"
  When call ./gru bios get --config "${GRU_CONF}" 127.0.0.1:5000
  The status should equal 0
  # check for some arbitrary keys to ensure it is not junk
  The stdout should include 'ProcessorHyperThreadingDisable'
  The stdout should include 'SRIOVEnable'
  The stdout should include 'VTdSupport'
  The lines of stdout should equal 866
  The lines of stderr should equal 1
End

# TODO: restore when args work with piping
# # neglecting to add a host as an arg should fail with instructions
# It "--config ${GRU_CONF}"
#   When call ./gru --config "${GRU_CONF}" bios get
#   The status should equal 1
#   The stderr should include 'Error: requires at least 1 arg(s), only received 0'
# End

# getting pending changes should return an error if the Bios/Settings.Attributes does not exist
It "--config ${GRU_CONF} --pending 127.0.0.1:5000"
  When call ./gru bios get --config "${GRU_CONF}" --pending 127.0.0.1:5000
  The status should equal 0
  The stdout should include 'Error'
  The stdout should include '"Attributes" does not exist or is null, the BIOS/firmware may need to updated for proper Attributes support'
  The lines of stdout should equal 2
  The lines of stderr should equal 1
End

# TODO: restore when marshaling JSON errors is fixed.
# getting pending changes should return an error if the Bios/Settings.Attributes does not exist and be valid json
#It "--config ${GRU_CONF} --pending 127.0.0.1:5000 --json"
#  When call ./gru bios get --config "${GRU_CONF}" --pending 127.0.0.1:5000 --json
#  The status should equal 0
#  The stdout should include 'error'
#  The stdout should include '\"Attributes\" does not exist or is null, the BIOS/firmware may need to updated for proper Attributes support'
#  The stdout should be_json
#  The lines of stderr should equal 0
#End

# getting keys from a file should return those keys
It "--config ${GRU_CONF} --from-file ${GRU_BIOS_KV} 127.0.0.1:5000"
  When call ./gru bios get --config "${GRU_CONF}" --from-file "${GRU_BIOS_KV}" 127.0.0.1:5000
  The status should equal 0
  The stdout should include 'BootTimeout'
  The stdout should include 'SRIOVEnable'
  The lines of stdout should equal 4
  The lines of stderr should equal 1
End

# getting keys from a file should return those keys and be valid json
It "--config ${GRU_CONF} --from-file ${GRU_BIOS_KV} 127.0.0.1:5000 --json"
  When call ./gru bios get --config "${GRU_CONF}" --from-file "${GRU_BIOS_KV}" 127.0.0.1:5000 --json
  The status should equal 0
  The stdout should include 'BootTimeout'
  The stdout should include 'SRIOVEnable'
  The stdout should be_json
  The lines of stderr should equal 0
End

# passing a shortcut should return a limited set of pre-defined keys
It "--config ${GRU_CONF} --virtualization 127.0.0.1:5003"
  When call ./gru bios get --config "${GRU_CONF}" --virtualization 127.0.0.1:5003
  The status should equal 0
  The stdout should include 'ProcAmdIOMMU'
  The stdout should include 'Sriov'
  The stdout should include 'ProcAmdVirtualization'
  The lines of stdout should equal 6
  The lines of stderr should equal 1
End

# passing a shortcut should return a limited set of pre-defined keys and be valid json
It "--config ${GRU_CONF} --virtualization 127.0.0.1:5003 --json"
  When call ./gru bios get --config "${GRU_CONF}" --virtualization 127.0.0.1:5003 --json
  The status should equal 0
  The stdout should include 'ProcAmdIOMMU'
  The stdout should include 'ProcAmdVirtualization'
  The stdout should include 'Sriov'
  The stdout should be_json
  The lines of stderr should equal 0
End

# piping in hosts should also work
Describe 'validate STDIN works'
  Data
    #|host1 127.0.0.1:5003
  End
  It "echo 127.0.0.1:5003 | --config ${GRU_CONF} --virtualization 127.0.0.1:5003 --json"
    When call ./gru bios get --config "${GRU_CONF}" --virtualization 127.0.0.1:5003 --json
    The status should equal 0
    The stdout should include 'ProcAmdIOMMU'
    The stdout should include 'ProcAmdVirtualization'
    The stdout should include 'Sriov'
    The stdout should be_json
    The lines of stderr should equal 0
  End
End

End
