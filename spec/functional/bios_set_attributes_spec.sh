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

Describe 'gru bios set'

BeforeAll use_valid_config

# FIXME: the simulator has a limitation and returns 405 for PATCH ops
# # setting specific keys should work
# It "--config ${GRU_CONF} --attributes ProcessorHyperThreadingDisable=0,SRIOVEnable=0 127.0.0.1:5000"
#   When call ./gru bios set --config "${GRU_CONF}" --attributes ProcessorHyperThreadingDisable=0,SRIOVEnable=0 127.0.0.1:5000
#   The status should equal 0
#   The stdout should include 'BIOS change(s) may be applied at: ["Immediate","OnReset","AtMaintenanceWindowStart","InMaintenanceWindowOnReset"]'
#   The lines of stdout should equal 3
# End

# TODO: restore when args work with piping
# # it should error if no attributes are passed to the flag
# It "--config ${GRU_CONF} --attributes 127.0.0.1:5000"
#   When call ./gru bios set --config "${GRU_CONF}" --attributes 127.0.0.1:5000
#   The status should equal 1
#   The stderr should include 'requires at least 1 arg(s), only received 0'
# End

End
