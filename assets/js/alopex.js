!(function () {
    if (window['_']) {
        return;
    }
    // 唯一字串
    window.uuid = function () {
        let part = function () {
            return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
        };
        return (part() + part() + '-' + part() + '-' + part() + '-' + part() + '-' + part() + part() + part());
    };
    // 日志重构
    (window.log = window.console.log) && (delete window.console);
    // 字串扩展
    String.prototype.trim = function (str) {
        str = str + '';
        str = (str !== undefined) || (str !== '') ? str : ' ';
        return this.replace(new RegExp('^(\\' + str + ')+|(\\' + str + ')+$', 'g'), '', '');
    };
    String.prototype.ltrim = function (str) {
        return this.replace(new RegExp('^(\\' + str + ')+', 'g'), '');
    };
    String.prototype.rtrim = function (str) {
        return this.replace(new RegExp('(\\' + str + ')+$', 'g'), '');
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
    // 文件上传地址
    window.common = {upload_url: 'http://127.0.0.1:81/backend/qp/common/upload/sdfsd'};
    // 认证用户数据
    window.auth_info = {
        token: '',
        random_str: '',
    }
    window.auth_info.token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNTYyNjg4MzMwIiwicmFuZG9tX3N0ciI6ImM2Yjg3YmNmN2M5M2M5OTBlZjRiMDk2ZWEwYjc1OTg5IiwidXNlcl9pZCI6IjEifQ==.Fy66uiZDihLTEJvO+NibJ1T7gpgIYdtvgoMEcpmN5Hw=';
    window.auth_info.random_str = 'c6b87bcf7c93c990ef4b096ea0b75989';
    // 列表筛选日期
    window.filter_daterange = {};
    // 列表页面数据
    window.list_data = {};
    // 弹框表单数据
    window.form_data = {};
    // 全局页面数据
    _.page = {
        logo_image: './imgs/logo.jpg',
        home_url: '/',
        userinfo_url: '#userinfo',
        refresh_url: '#refresh',
        logout_url: '#logout',
        clean_url: '#clean',
        copyright: '&copy;2019~2024 万里独行侠  田伯光',
    };
    //
    // 标签页
    $('#body-tabs').tabs({
        tabs: [{
            title: '表格',
            type: 'ajax',
            icon: 'icon-star',
            url: './table.html',
            forbidClose: true,
        }],
        onLoad: function (tab) {
            if (tab.url === './table.html') {
                for (let id in filter_daterange) {
                    $('#' + id).datetimepicker('remove');
                }
                list_data = {
                    url: 'http://127.0.0.1:81/backend/qp/admin/list',
                    cols: [
                        {
                            label: '编号',
                            name: 'id',
                            width: 'auto',
                            sort: true,
                        },
                        {
                            label: '用户名',
                            name: 'username',
                            width: 'auto',
                        },
                        {
                            label: '密码',
                            name: 'password',
                            width: 'auto',
                        },
                        {
                            label: '邮箱',
                            name: 'email',
                            width: 'auto',
                        },
                        {
                            label: '注册时间',
                            name: 'created_at',
                            width: 150,
                        },
                        {
                            label: '状态',
                            name: 'status',
                            width: 100,
                            html: true,
                        },
                    ],
                };
            }
        },
    });
    // 事件
    $('body').on('shown.zui.modal', '#triggerModal', function () {
        form_data['editor'] = form_data['uploader'] = form_data['daterange'] = {};
        let self = $(this);
        // 时间区间
        let daterange = self.find('.form-daterange');
        if (daterange.length > 0) {
            $.each(daterange, function (i, tag) {
                let uuidStr = 'daterange-' + uuid();
                $(tag).attr('id', uuidStr);
                form_data['daterange'][uuidStr] = $(tag).datetimepicker({
                    language: 'zh-CN',
                    weekStart: 1,
                    todayBtn: 1,
                    autoclose: 1,
                    todayHighlight: 1,
                    startView: 2,
                    minView: 2,
                    forceParse: 0,
                    format: 'yyyy-mm-dd',
                });
            });
        }
        // 上传
        let uploader = self.find('.uploader');
        if (uploader.length > 0) {
            $.each(uploader, function (i, tag) {
                let uuidStr = 'uploader-' + uuid();
                $(tag).attr('id', uuidStr);
                form_data['uploader'][uuidStr] = $(tag).uploader({
                    url: common.upload_url,
                    prevent_duplicates: true,
                    unique_names: true,
                    autoUpload: true,
                    runtimes: 'html5',
                    max_retries: 0,
                    rename: true,
                    headers: {
                        'Debug': 'true',
                        'Token': auth_info.token,
                        'Random_str': auth_info.random_str,
                    },
                    onUploadComplete: function (files) {
                        let otag = $(this.$[0]).parents('.form-group');
                        let value = otag.data('value');
                        value = $.isArray(value) ? value : [];
                        $.each(files, function (k, v) {
                            if (v.url) {
                                value[value.length] = v.url;
                            }
                        });
                        otag.data('value', value);
                    },
                });
            });
        }
        // 富文本
        let editor = self.find('.form-editor');
        if (editor.length > 0) {
            let E = window.wangEditor;
            $.each(editor, function (i, tag) {
                let uuidStr = 'editor-' + uuid();
                $(tag).attr('id', uuidStr);
                let editor = new E('#' + uuidStr);
                editor.customConfig.uploadImgShowBase64 = true;
                editor.create();
                form_data['editor'][uuidStr] = editor;
            });
        }
        // 下拉框
        let selecter = self.find('.chosen-select');
        if (selecter.length > 0) {
            $.each(selecter, function (i, tag) {
                $(tag).chosen({no_results_text: '没有找到'});
            });
        }
        // 单选/复选框
        let checker = self.find('.checkbox-primary,.radio-primary')
        if (checker.length > 0) {
            $.each(checker, function (i, tag) {
                let [input, label] = [$(tag).find('input'), $(tag).find('label')];
                if (input && label) {
                    let uuidStr = 'k-' + uuid();
                    input.attr('id', uuidStr);
                    label.attr('for', uuidStr);
                }
            });
        }
    }).on('hidden.zui.modal', '#triggerModal', function () {
        for (let id in form_data['daterange']) {
            $('#' + id).datetimepicker('remove');
        }
        form_data['editor'] = form_data['uploader'] = form_data['daterange'] = {};
    }).on('click', '#body-modal .has-error', function () {
        $(this).removeClass('has-error').find('.help-block').html('');
    }).on('change', '#body-modal .form-daterange', function () {
        let parent = $(this).parent().parent();
        let tags = parent.find('.form-daterange');
        let value = [];
        if (tags.length > 0) {
            $.each(tags, function (i, tag) {
                value[value.length] = ($(tag).val() + '').trim();
            });
        }
        parent.data('value', value);
    }).on('change', '#body-modal .switch>[type=checkbox]', function () {
        $(this).parents('.form-group').data('value', $(this).is(':checked') ? '1' : '0');
    }).on('change', '#body-modal select.chosen-select', function () {
        $(this).parents('.form-group').data('value', ($(this).val() + '').trim());
    }).on('change', '#body-modal .checkbox-primary>[type=checkbox]', function () {
        let parent = $(this).parent().parent().parent();
        let tags = parent.find('[type=checkbox]:checked');
        let value = [];
        if (tags.length > 0) {
            $.each(tags, function (i, tag) {
                value[value.length] = ($(tag).val() + '').trim();
            });
        }
        parent.data('value', value);
    }).on('change', '#body-modal .radio-primary>[type=radio]', function () {
        let parent = $(this).parent().parent().parent();
        let tags = parent.find('[type=radio]:checked');
        let value = '';
        if (tags.length > 0) {
            value = ($(tags[0]).val() + '').trim();
        }
        parent.data('value', value);
    }).on('click', '#body-modal .form-submit', function () {
        let form = $(this).parent().parent();
        let targetId = form.data('id');
        if (('' !== targetId) && (!/^[1-9][0-9]*$/.test(targetId))) {
            new $.zui.Messager('操作配置异常！', {
                icon: 'warning-sign',
                placement: 'center',
                type: 'danger',
                close: true,
                actions: [{
                    name: 'refresh',
                    icon: 'refresh',
                    text: '刷新',
                    action: function () {
                        location.reload();
                    }
                }]
            }).show();
            return;
        }
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
                    } else if ('editor' === type) {
                        let id = otag.find('.form-editor').attr('id');
                        let tmp = form_data['editor'][id];
                        if (undefined !== tmp) {
                            value = tmp.txt.html();
                            if ('<p><br></p>' === value) {
                                value = '';
                            }
                        }
                    } else if ('range' === type) {
                        let inputs = otag.find('input.form-control');
                        if (inputs.length === 2) {
                            let [a, b] = [$(inputs[0]).val().trim(), $(inputs[1]).val().trim()];
                            if (a || b) {
                                value = [a, b];
                            }
                        }
                    } else {
                        value = otag.data('value');
                        if (undefined === value) {
                            value = '';
                        }
                    }
                    if (must && (('' === value) || ($.isArray(value) && (value.length === 0)))) {
                        otag.addClass('has-error').append('<p class=\'help-block\'>' + label + '不能为空！</p>');
                        isOK = false;
                        return;
                    }
                    if ((typeof value === 'string') && ('' !== regex)) {
                        let regs = regex.split(' ');
                        let reg = regs.length > 1 ? (new RegExp(regs[0], regs[1])) : (new RegExp(regs[0]));
                        if (!reg.test(value)) {
                            otag.addClass('has-error').append('<p class=\'help-block\'>' + label + '校验失败！</p>');
                            isOK = false;
                            return;
                        }
                    }
                    result[name] = value;
                }
            });
            if (!isOK) {
                return;
            }
            if (targetId)
                form.data('data', result);
        }
    });
});