$(function () {
    'use strict';

    var
        id,
        ip,
        idCheck = false,
        ipCheck = false,
        $id = $('#CreateMachineId'),
        $ip = $('#CreateMachineIp'),
        idInputErrorFunc = function () {
            FormValidatorUI.inputError($id, "input id error", "the id range 1 ~ 1024");
        },
        ipInputErrorFunc = function () {
            FormValidatorUI.inputError($ip, "input ip error", "the ip format is incorrect");
        };

    $id.on('change', function () {
        id = $id.val();
        if (id > 1 << 10 || id < 1) {
            idInputErrorFunc();
            idCheck = false;
        } else {
            FormValidatorUI.inputSuccess($id, "input id success", "");
            idCheck = true;
        }
    });

    $ip.on('change', function () {
        ip = $ip.val();
        var re = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
        if (!re.test(ip)) {
            ipInputErrorFunc();
            ipCheck = false;
        } else {
            FormValidatorUI.inputSuccess($ip, "input ip success", "");
            ipCheck = true;
        }
    });

    $("#CreateMachineSubmitButton").on('click', function () {
        if (!idCheck) {
            idInputErrorFunc();
            return false;
        }

        if (!ipCheck) {
            ipInputErrorFunc();
            return false;
        }

        $.post('/machine/store', {
            "id": id,
            "ip": ip,
        }, function (data, status) {
            console.log(data, status);
        });

    });

    $('#MachineListTable').delegate('.MachineItemDelBtn', 'click', function () {
        var $this = $(this);
        var $btn = $(this).button('loading');
        var ip = $this.attr('data-ip');
        if (!ip) {
            $btn.button('reset');
            return false
        }

        $.post('/machine/delete', {
            'ip': ip,
        }, function (data, status) {
            $btn.button('reset');
            console.log(data)
            if (data.success) {
                toastr.success('Delete '+ ip + ' success!')
            } else {
                toastr.error(data.message)
            }
        });
    });

});