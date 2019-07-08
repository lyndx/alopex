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
    // 认证数据
    window.authinfo = {token: '', random_str: ''}
    authinfo.token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNTYyNTc0MDg2IiwicmFuZG9tX3N0ciI6Ijg5MjQ0OGU0OTIxMDRhYzNmYmZjMThjMjUyZTg1ODdhIiwidXNlcl9pZCI6IjEifQ==.E1ZtYMEBIUS26TKNaxSXWJMJPGleG8DAdWcDRfzIu+8=';
    authinfo.random_str = '892448e492104ac3fbfc18c252e8587a';
    // 表单数据
    window.fdata = {};
    // 列表数据
    window.table = {};
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
    // 标签页
    $('#body-tabs').tabs({
        tabs: [{
            forbidClose: true,
            title: '表格',
            icon: 'icon-star',
            url: './table.html',
            type: 'ajax'
        }],
        onLoad: function (tab) {
            if (tab.url == './table.html') {
                for (let id in table['daterange']) {
                    $('#' + id).datetimepicker('remove');
                }
            }
        }
    });
    // 事件
    $('body').on('shown.zui.modal', '#triggerModal', function () {
        let [uploader, selecter, editor, daterange] = [$(this).find('.uploader'), $(this).find('.chosen-select'), $(this).find('.form-editor'), $(this).find('.form-daterange')];
        fdata['editor'] = fdata['uploader'] = fdata['daterange'] = {};
        // 文件上传
        if (uploader.length > 0) {
            $.each(uploader, function (i, tag) {
                let id = $(tag).attr('id');
                if (id) {
                    fdata['uploader'][id] = $(tag).uploader({
                        url: 'http://127.0.0.1:81/backend/qp/common/upload/sdfsd',
                        prevent_duplicates: true,
                        unique_names: true,
                        autoUpload: true,
                        runtimes: 'html5',
                        max_retries: 0,
                        rename: true,
                        headers: {
                            'Debug': 'true',
                            'Token': authinfo.token,
                            'Random_str': authinfo.random_str
                        },
                        onUploadComplete: function (files) {
                            let otag = $(this.$[0]).parent();
                            let value = otag.data('value');
                            value = $.isArray(value) ? value : [];
                            $.each(files, function (k, v) {
                                if (v.url) {
                                    value[value.length] = v.url;
                                }
                            });
                            otag.data('value', value);
                        }
                    });
                }
            });
        }
        // 选择框
        if (selecter.length > 0) {
            $.each(selecter, function (i, tag) {
                $(tag).chosen({no_results_text: '没有找到'});
            });
        }
        // 富文本
        if (editor.length > 0) {
            var E = window.wangEditor;
            $.each(editor, function (i, tag) {
                let id = $(tag).attr('id');
                if (id) {
                    let editor = new E('#' + id);
                    editor.customConfig.uploadImgShowBase64 = true;
                    editor.create()
                    fdata['editor'][id] = editor;
                }
            });
        }
        // 时间区间
        if (daterange.length > 0) {
            $.each(daterange, function (i, tag) {
                let id = $(tag).attr('id');
                if (id) {
                    fdata['daterange'][id] = $(tag).datetimepicker({
                        language: "zh-CN",
                        weekStart: 1,
                        todayBtn: 1,
                        autoclose: 1,
                        todayHighlight: 1,
                        startView: 2,
                        minView: 2,
                        forceParse: 0,
                        format: "yyyy-mm-dd"
                    });
                }
            });
        }
    }).on('hidden.zui.modal', '#triggerModal', function () {
        for (let id in fdata['daterange']) {
            $('#' + id).datetimepicker('remove');
        }
        fdata['editor'] = fdata['uploader'] = fdata['daterange'] = {};
    }).on('click', '#body-modal .has-error', function () {
        $(this).removeClass('has-error').find('.help-block').html('');
    }).on('change', '#body-modal select.chosen-select', function () {
        let value = $(this).val();
        $(this).parents('.form-group').data('value', value);
    }).on('click', '#body-modal .form-submit', function () {
        let form = $(this).parent().parent();
        if (form.find('.has-error').length > 0) {
            form.data('data', {});
        }
        let [fields, isOK, result] = [form.find('.form-group'), true, {}];
        if (fields.length > 0) {
            $.each(fields, function (i, tag) {
                let otag = $(tag);
                otag.find('.help-block').remove();
                let [name, label, must, type, regex, value] = [otag.attr('f-name'), otag.find('label').text(), otag.attr('f-must'), otag.attr('f-type'), otag.attr('f-regex'), ''];
                name = (undefined !== name) && ('' !== name.trim()) ? name.trim() : '';
                if ('' !== name) {
                    must = (undefined !== must) && ('true' === must.trim().toLowerCase());
                    type = (undefined !== type) && ('' !== type.trim().toLowerCase()) ? type.trim().toLowerCase() : 'string';
                    label = (undefined !== label) ? label.trim() : name;
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
                            if (start || end) {
                                value = [start, end];
                            }
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
                    } else if ('editor' === type) {
                        let id = otag.find('.form-editor').attr('id');
                        let tmp = fdata['editor'][id]
                        if (undefined !== tmp) {
                            value = tmp.txt.html();
                            if ('<p><br></p>' === value){
                                value = '';
                            }
                        }
                        log(value)
                    } else {
                        value = otag.data('value');
                        if (undefined === value) {
                            value = '';
                        }
                    }

                    if (must && (('' === value) || ($.isArray(value) && (value.length === 0)))) {
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
            form.data('data', result);
        }
    });
});