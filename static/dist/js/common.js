(function ($, exports) {
    var FormValidatorUI = function () {
        this._input = function (el, title, help, level) {
            var fa, divClass;
            switch (level) {
                case 'success':
                    fa = 'fa-check';
                    divClass = 'has-success';
                    break;
                case 'warning':
                    fa = 'fa-bell-o';
                    divClass = 'has-warning';
                    break;
                case 'error':
                    fa = 'fa-times-circle-o';
                    divClass = 'has-error';
                    break;
                default:
                    fa = 'fa-bell-o';
                    divClass = 'has-warning';
                    break;
            }
            var $elPrev = el.prev();
            $elPrev.html($('<i class="fa ' + fa +'"></i> '));
            $elPrev.append(' ' + title);
            if (help) {
                if (el.next('span').length <= 0) {
                    el.after($('<span class="help-block">'+help+'</span>'));
                }
            } else {
                el.next('span').remove();
            }

            el.parent().attr('class', 'form-group ' + divClass);
            return;
        };

        this.inputSuccess =  function (el, title, help) {
            this._input(el, title, help, 'success');
        };

        this.inputWarning =  function (el, title, help) {
            this._input(el, title, help, 'warning');
        };

        this.inputError =  function (el, title, help) {
            this._input(el, title, help, 'error');
        };
    };
    exports.FormValidatorUI = new FormValidatorUI();
})(jQuery, window);

