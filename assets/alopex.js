/* Alopex.js */
!(function () {
    if (window['_']) return;
    // 日志重构
    window.log = window.console.log
    delete window.console;
    // 字串扩展
    String.prototype.endWith = function (endStr) {
        var d = this.length - endStr.length;
        return (d >= 0 && this.lastIndexOf(endStr) == d);
    };
    String.prototype.beginWith = function (endStr) {
        return this.indexOf(endStr) == 0;
    };
    // 数据绑定
    window._ = new Proxy({}, {
        set: function (obj, prop, value) {
            obj[prop] = value;
            // 数据
            $.each($('[x-data=' + prop + ']'), function (i, tag) {
                // 属性
                $.each(tag.attributes, function (ii, attr) {
                    let [nodeName, nodeValue] = [$.trim(attr.nodeName), $.trim(attr.value)];
                    if ((nodeName != 'x-data') && nodeName.beginWith('x-')) {
                        nodeName = nodeName.replace('x-', '');
                        nodeValue = obj[prop][nodeValue] != undefined ? obj[prop][nodeValue] : $(tag).attr(nodeName);
                        $(tag).attr(nodeName, nodeValue);
                    }
                });
                // 内容
                let htmlStr = $(tag).html();
                let fields = htmlStr.match(/\{\{[^\{\}\r\n]+\}\}/g);
                $.each(fields, function (ii, field) {
                    let trimField = $.trim(field.replace('{{', '').replace('}}', ''));
                    if (trimField.beginWith(prop + '.')) {
                        trimField = trimField.replace(prop + '.', '');
                        let value = obj[prop][trimField] ? obj[prop][trimField] : '';
                        htmlStr = htmlStr.replace(field, value);
                    } else if (trimField == prop) {
                        htmlStr = htmlStr.replace(field, obj[prop]);
                    }
                });
                $(tag).html(htmlStr);
            });
            let field = obj[prop] != undefined ? obj[prop] : false;
            // 列表
            $.each($('[x-list=' + prop + ']'), function (i, tag) {
                $(tag).html('');
                let items = field ? field : [];
                if (items.length > 0) {
                    let template = $.trim($(tag).attr('x-template'));
                    let fields = template.match(/\{\{[^\{\}\r\n]+\}\}/g);
                    let htmlStr = [];
                    $.each(items, function (ii, item) {
                        let itemStr = template;
                        $.each(fields, function (iii, field) {
                            let trimField = $.trim(field.replace('{{', '').replace('}}', ''));
                            let value = item[trimField] ? item[trimField] : '';
                            itemStr = itemStr.replace(field, value);
                        });
                        htmlStr[htmlStr.length] = itemStr;
                    });
                    $(tag).html(htmlStr.join(''));
                }
            });
            // 切换
            $.each($('[x-toggle]'), function (i, tag) {
                if ($.trim($(tag).attr('x-toggle')) == prop) {
                    let nodeClass = $.trim($(tag).attr('x-class'));
                    if (field) {
                        if (nodeClass != "") {
                            $(tag).addClass(nodeClass)
                        } else {
                            $(tag).show();
                        }
                    } else {
                        if (nodeClass != "") {
                            $(tag).removeClass(nodeClass)
                        } else {
                            $(tag).hide();
                        }
                    }
                }
            });
            // HTML片段
            $.each($('[x-html=' + prop + ']'), function (i, tag) {
                let htmlPath = field ? (field + '') : '';
                if (htmlPath == '') {
                    $(tag).html('');
                }
                if (!htmlPath.endWith(".html")) {
                    htmlPath += '.html';
                }
                $.ajax({
                    url: htmlPath.toLowerCase(),
                    cache: true,
                    success: function (htmlStr) {
                        $(tag).html(htmlStr);
                        _.dd = {xc: "sdfsdfsdfsdferd"}
                    },
                    error: function () {
                        $(tag).html('');
                    }
                });
            });
            return true;
        }
    });
    // 焦点时间
    $('body').on('focus', 'input', function () {
        let parent = $(this).parent();
        parent.removeClass('error');
        parent.find('.tip').text('');
    }).on('blur', 'input', function () {
        let v = $(this).val();
        _.fd = {value: v};
        if (v == "") {
            let parent = $(this).parent();
            parent.addClass('error');
            parent.find('.tip').text('书屋错误');
        }
    });
    // 模拟数据
    _.ss = {value: "sdsdsd", placeholder: "sdfsdfsd"};
    _.aa = [
        {name: 'sdfs', value: 234242},
        {name: 'sdfs', value: 234242},
        {name: 'sdfs', value: 234242},
        {name: 'sdfs', value: 234242},
        {name: 'sdfs', value: 234242},
    ];
    _.ww = 23;
    _.xx = 'X'
    _.a = {x: "axxxxxxxxx"}
    _.fd = {value: 234234, placeholder: '请输入验证码'};
    _.title = "Alopex管理端"
})();




