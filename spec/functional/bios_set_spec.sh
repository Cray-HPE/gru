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


Describe 'gru set bios'

BeforeAll use_valid_config
BeforeAll use_valid_bios_attributes_file

# setting without flags should fail
It "--config ${GRU_CONF} 127.0.0.1:5000"
  When call ./gru set bios --config "${GRU_CONF}" 127.0.0.1:5000
  The status should equal 1
  The stderr should include 'An error occurred: at least one of the flags in the group [attributes from-file virt defaults] is required'
End

# neglecting to add a host as an arg should fail with instructions
It "--config ${GRU_CONF}"
  When call ./gru --config "${GRU_CONF}" set bios
  The status should equal 1
  The stderr should include 'Error: requires at least 1 arg(s), only received 0'
End

# FIXME: the emulator does not support PATCH ops.  Once it does, uncomment, adjust if needed
# # restoring defaults
# It "--config ${GRU_CONF} --defaults 127.0.0.1:5000"
#   When call ./gru set bios --config "${GRU_CONF}" --defaults 127.0.0.1:5000
#   The status should equal 0
#   The stdout should include 'Error'
#   The stdout should include 'BIOS reset failure: unable to execute request, no target provided'
#   The lines of stdout should equal 3
# End

# # restoring defaults and be valid json
# It "--config ${GRU_CONF} --defaults 127.0.0.1:5000 --json"
#   When call ./gru set bios --config "${GRU_CONF}" --defaults 127.0.0.1:5000 --json
#   The status should equal 0
#   The stdout should include 'Error'
#   The stdout should include 'BIOS reset failure: unable to execute request, no target provided'
#   The stdout should be_json
# End

# # setting keys from a file should return those keys
# It "--config ${GRU_CONF} --from-file ${GRU_BIOS_KV} 127.0.0.1:5000"
#   When call ./gru set bios --config "${GRU_CONF}" --from-file "${GRU_BIOS_KV}" 127.0.0.1:5000
#   The status should equal 0
#   The stdout should include 'Pending'
#   The stdout should include 'Attributes'
#   The stdout should include 'BootMode'
#   The stdout should include 'BootTimeout'
#   The stdout should include 'ProcessorHyperThreadingDisable'
#   The stdout should include 'ProcessorVmxEnable'
#   The stdout should include 'ProcessorX2apic'
#   The stdout should include 'SRIOVEnable'
#   The stdout should include 'SvrMngmntAcpiIpmi'
#   The stdout should include 'VTdSupport'
#   The lines of stdout should equal 3
# End

# # setting keys from a file should return those keys and be valid json
# It "--config ${GRU_CONF} --from-file ${GRU_BIOS_KV} 127.0.0.1:5000 --json"
#   When call ./gru set bios --config "${GRU_CONF}" --from-file "${GRU_BIOS_KV}" 127.0.0.1:5000 --json
#   The status should equal 0
#   The stdout should include 'Pending'
#   The stdout should include 'Attributes'
#   The stdout should include 'BootMode'
#   The stdout should include 'BootTimeout'
#   The stdout should include 'ProcessorHyperThreadingDisable'
#   The stdout should include 'ProcessorVmxEnable'
#   The stdout should include 'ProcessorX2apic'
#   The stdout should include 'SRIOVEnable'
#   The stdout should include 'SvrMngmntAcpiIpmi'
#   The stdout should include 'VTdSupport'
#   The stdout should be_json
# End

# # passing a shortcut should return a limited set of pre-defined keys
# It "--config ${GRU_CONF} --virt 127.0.0.1:5000"
#   When call ./gru set bios --config "${GRU_CONF}" --virt 127.0.0.1:5000
#   The status should equal 0
#   The stdout should include 'BootMode'
#   The stdout should include 'ProcessorHyperThreadingDisable'
#   The stdout should include 'ProcessorVmxEnable'
#   The stdout should include 'ProcessorX2apic'
#   The stdout should include 'SRIOVEnable'
#   The stdout should include 'SvrMngmntAcpiIpmi'
#   The stdout should include 'VTdSupport'
# End

# # passing a shortcut should return a limited set of pre-defined keys and be valid json
# It "--config ${GRU_CONF} --virt 127.0.0.1:5000 --json"
#   When call ./gru set bios --config "${GRU_CONF}" --virt 127.0.0.1:5000 --json
#   The status should equal 0
#   The stdout should include 'BootMode'
#   The stdout should include 'ProcessorHyperThreadingDisable'
#   The stdout should include 'ProcessorVmxEnable'
#   The stdout should include 'ProcessorX2apic'
#   The stdout should include 'SRIOVEnable'
#   The stdout should include 'SvrMngmntAcpiIpmi'
#   The stdout should include 'VTdSupport'
#   The stdout should be_json
# End

End
