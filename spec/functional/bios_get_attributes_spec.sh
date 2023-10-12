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

Describe 'gru get bios'

BeforeAll use_valid_config

# getting a single key should return only those key
It "--config ${GRU_CONF} --attributes BootTimeout 127.0.0.1:5000"
  When call ./gru get bios --config "${GRU_CONF}" --attributes BootTimeout 127.0.0.1:5000
  The status should equal 0
  The stdout should include 'BootTimeout'
  The lines of stdout should equal 3
End

# getting a single key should return only those key in json
It "--config ${GRU_CONF} --attributes BootTimeout 127.0.0.1:5000 --json"
  When call ./gru get bios --config "${GRU_CONF}" --attributes BootTimeout 127.0.0.1:5000 --json
  The status should equal 0
  The stdout should include 'BootTimeout'
  The stdout should be_json
End

# getting multiple keys should return only those keys
It "--config ${GRU_CONF} --attributes ProcessorHyperThreadingDisable,SRIOVEnable 127.0.0.1:5000"
  When call ./gru get bios --config "${GRU_CONF}" --attributes ProcessorHyperThreadingDisable,SRIOVEnable 127.0.0.1:5000
  The status should equal 0
  The stdout should include 'ProcessorHyperThreadingDisable'
  The stdout should include 'SRIOVEnable'
  The lines of stdout should equal 4
End

# getting specific keys should return only those keys and should be json
It "--config ${GRU_CONF} --attributes ProcessorHyperThreadingDisable,SRIOVEnable 127.0.0.1:5000 --json"
  When call ./gru get bios --config "${GRU_CONF}" --attributes ProcessorHyperThreadingDisable,SRIOVEnable 127.0.0.1:5000 --json
  The status should equal 0
  The stdout should include 'ProcessorHyperThreadingDisable'
  The stdout should include 'SRIOVEnable'
  The stdout should be_json
End

# it should error if no matching keys were found
It "--config ${GRU_CONF} --attributes junk 127.0.0.1:5000"
  When call ./gru get bios --config "${GRU_CONF}" --attributes junk 127.0.0.1:5000
  The status should equal 0
  The line 3 of stdout should include 'junk'
  The line 3 of stdout should include '<nil>'
  The lines of stdout should equal 3
End

# it should error if no attributes are passed to the flag
It "--config ${GRU_CONF} --attributes 127.0.0.1:5000"
  When call ./gru get bios --config "${GRU_CONF}" --attributes 127.0.0.1:5000
  The status should equal 1
  The stderr should include 'requires at least 1 arg(s), only received 0'
End

# --virt shortcut should return only virtualization attributes (Gigabyte)
It "--config ${GRU_CONF} --virt 127.0.0.1:5001"
  When call ./gru get bios --config "${GRU_CONF}" --virt 127.0.0.1:5001
  The status should equal 0
  The stdout should include 'SR-IOV Support'
  The stdout should include 'SMT Control'
  The stdout should include 'Local APIC Mode'
  The stdout should include 'IOMMU'
  The stdout should include 'SVM Mode'
  The lines of stdout should equal 7
End

# --virt shortcut should return only virtualization attributes in json format (Gigabyte)
It "--config ${GRU_CONF} --virt 127.0.0.1:5001 --json"
  When call ./gru get bios --config "${GRU_CONF}" --virt 127.0.0.1:5001 --json
  The status should equal 0
  The stdout should include 'PCIS007' # 'SR-IOV Support'
  The stdout should include 'Rome0039' # 'Local APIC Mode'
  The stdout should include 'Rome0059' # 'SMT Control'
  The stdout should include 'Rome0162' # 'IOMMU'
  The stdout should include 'Rome0565' # 'SVM Mode'
  The stdout should be_json
End


End
