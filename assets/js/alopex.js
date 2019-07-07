!(function () {
    if (window['_']) {
        return;
    }
    // 日志重构
    (window.log = window.console.log) && (delete window.console);
    // 字串扩展
    String.prototype.trim = function (str) {
        str = str + '';
        str = (str !== undefined) || (str !== '') ? str : ' ';
        return this.replace(new RegExp('^\\' + str + '+|\\' + str + '+$', 'g'), '', '');
    };
    String.prototype.ltrim = function (str) {
        return this.replace(new RegExp('^\\' + str + '+', 'g'), '');
    };
    String.prototype.rtrim = function (str) {
        return this.replace(new RegExp('\\' + str + '+$', 'g'), '');
    };
    String.prototype.hasPrefix = function (str) {
        return this.indexOf(str) === 0;
    };
    String.prototype.hasSuffix = function (str) {
        let d = this.length - str.length;
        return (d >= 0 && this.lastIndexOf(str) === d);
    };
    // 数据绑定
    window._ = new Proxy({}, {
        set: function (obj, prop, value) {
            if ((undefined !== obj[prop]) && (obj[prop] === value)) {
                return false;
            }
            // 赋值
            obj[prop] = value;
            // 数据
            let [pvalue, t_attrs, t_list, t_toggle, t_html] = [(undefined !== obj[prop]) ? obj[prop] : '', $('[x-data=' + prop + ']'), $('[x-list=' + prop + ']'), $('[x-toggle]'), $('[x-html=' + prop + ']')]
            if (t_attrs.length > 0) {
                $.each(t_attrs, function (i, tag) {
                    let [tobj, tattrs] = [$(tag), tag.attributes];
                    // 属性赋值
                    if (tattrs.length > 1) {
                        $.each(tattrs, function (ii, v) {
                            let [name, field] = [v.nodeName.trim(), v.value.trim()];
                            if (('x-data' !== name) && name.hasPrefix('x-')) {
                                name = name.ltrim('x-');
                                let ovalue = tobj.attr(name);
                                value = $.isPlainObject(pvalue) && (undefined !== pvalue[field]) ? pvalue[field] : (ovalue ? ovalue : '');
                                tobj.attr(name, value);
                            }
                        });
                    }
                    // 内容赋值
                    let hstr = tobj.html();
                    let fields = hstr.match(/\{\{[^\{\}\r\n]+\}\}/g);
                    if ($.isArray(fields) && (fields.length > 0)) {
                        $.each(fields, function (ii, v) {
                            let [field, value] = [v.ltrim('{{', '').rtrim('}}', '').trim(), ''];
                            if (prop === field) {
                                value = pvalue;
                            } else if (field.hasPrefix(prop + '.')) {
                                field = field.ltrim(prop + '.', '');
                                value = $.isPlainObject(pvalue) && (undefined !== pvalue[field]) ? pvalue[field] : '';
                            }
                            hstr = hstr.replace(v, value);
                        });
                        tobj.html(hstr);
                    }
                });
            }
            // 列表
            if (t_list.length > 0) {
                let list = pvalue.isArray() ? pvalue : [];
                $.each(t_list, function (i, tag) {
                    let otag = $(tag);
                    otag.html('');
                    if (list.length > 0) {
                        let template = otag.attr('x-template').trim();
                        let [fields, hstr] = [template.match(/\{\{[^\{\}\r\n]+\}\}/g), []];
                        $.each(list, function (ii, item) {
                            let istr = template;
                            if ($.isArray(fields) && (fields.length > 0)) {
                                $.each(fields, function (iii, v) {
                                    let field = v.ltrim('{{').rtrim('}}').trim();
                                    let value = $.isPlainObject(item) && (undefined !== item[field]) ? item[field] : '';
                                    istr = istr.replace(v, value);
                                });
                            }
                            hstr[hstr.length] = istr;
                        });
                        otag.html(hstr.join(''));
                    }
                });
            }
            // 切换
            if (t_toggle.length > 0) {
                $.each(t_toggle, function (i, tag) {
                    if (prop === $(tag).attr('x-toggle').trim()) {
                        let nclass = $(tag).attr('x-class').trim();
                        if (pvalue) {
                            if ('' !== nclass) {
                                $(tag).addClass(nclass)
                            } else {
                                $(tag).show();
                            }
                        } else {
                            if ('' !== nclass) {
                                $(tag).removeClass(nclass)
                            } else {
                                $(tag).hide();
                            }
                        }
                    }
                });
            }
            // HTML片段
            if (t_html.length > 0) {
                $.each($('[x-html=' + prop + ']'), function (i, tag) {
                    if ((typeof pvalue != 'string') || ('' === pvalue)) {
                        $(tag).html('');
                    }
                    if (!pvalue.hasSuffix('.html')) {
                        pvalue += '.html';
                    }
                    $.ajax({
                        url: pvalue.toLowerCase(),
                        cache: true,
                        error: function () {
                            $(tag).html('');
                        },
                        success: function (hstr) {
                            $(tag).html(hstr);
                        },
                    });
                });
            }
            return true;
        },
    });
})();
$(document).ready(function () {
    _.page_header = {
        logo_image: './assets/imgs/logo.jpg',
        home_url: '/',
        userinfo_url: '#userinfo',
        refresh_url: '#refresh',
        logout_url: '#logout',
        clean_url: '#clean',
    };
    _.page_footer = {copyright: '&copy;2019~2024 万里独行侠  田伯光'};
    // 焦点事件
    $('body').on('click', '.has-error', function () {
        $(this).removeClass('has-error').find('.help-block').html('');
    }).on('click', '.form-submit', function () {
        let form = $(this).parent().parent();
        if (form.find('.has-error').length > 0) {
            //return;
        }
        let [fields, isOK, result] = [form.find('.form-group'), true, {}];
        if (fields.length > 0) {
            $.each(fields, function (i, tag) {
                let otag = $(tag);
                otag.find('.help-block').remove();
                let [name, label, must, type, regex, value] = [otag.attr('f-name'), otag.find('label').text(), otag.attr('f-must'), otag.attr('f-type'), otag.attr('f-regex'), ''];
                name = (undefined !== name) && ('' !== name.trim()) ? name.trim() : '';
                if ('' !== name) {
                    label = (undefined !== label) ? label.trim() : name;
                    must = (undefined !== must) && ('true' === must.trim().toLowerCase());
                    type = (undefined !== type) && ('' !== type.trim().toLowerCase()) ? type.trim().toLowerCase() : 'string';
                    regex = (undefined !== regex) ? regex.trim() : '';
                    if (('input' === type) || ('textarea' === type)) {
                        value = otag.find(type).val().trim();
                    } else if ('daterange' === type) {
                        let tags = otag.find('input');
                        value = [];
                        if (tags.length === 2) {
                            let [start, end] = [$(tags[0]).val().trim(), $(tags[1]).val().trim()];
                            if (('' !== start) && ('' !== end) && (start > end)) {
                                $(tags[0]).val('');
                                start = '';
                            }
                            value = [start, end];
                        }
                    } else if ('switch' === type) {
                        let tags = otag.find('input:checked');
                        value = 0;
                        if (tags.length > 0) {
                            value = 1;
                        }
                    } else if ('checkbox' === type) {
                        let tags = otag.find('input:checked');
                        value = [];
                        if (tags.length > 0) {
                            $.each(tags, function (k, v) {
                                value[value.length] = $(v).val().trim();
                            });
                        }
                    } else if ('radio' === type) {
                        let tmp = otag.find('input:checked');
                        if (undefined !== tmp) {
                            value = tmp.val().trim();
                        }
                    } else {
                        value = otag.data('value');
                        if (undefined === value) {
                            value = '';
                        }
                    }
                    if (must && (('' === value) || ($.isArray(value) && (value.length < 1)))) {
                        otag.addClass('has-error').append('<p class="help-block">' + label + '不能为空！</p>');
                        isOK = false;
                        return;
                    }
                    if ((typeof value == 'string') && ('' !== regex)) {
                        let regs = regex.split(' ');
                        let reg = regs.length > 1 ? (new RegExp(regs[0], regs[1])) : (new RegExp(regs[0]));
                        if (!reg.test(value)) {
                            otag.addClass('has-error').append('<p class="help-block">' + label + '校验失败！</p>');
                            isOK = false;
                            return;
                        }
                    }
                    result[name] = value;
                }
            });
            if (!isOK) {
                result = {};
            }
            log(result);
        }
    });
});