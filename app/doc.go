package app

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"
)

func init() {
	Mux.HandleFunc("/doc", func(rep http.ResponseWriter, req *http.Request) {
		result := make(map[string]interface{})
		for _, module := range String("route").Scan("", true) {
			result[module] = (*route).GValue(module, false).ToMS()
		}
		bs, _ := json.Marshal(result)
		html := `<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="utf-8">
		    <title>接口清单</title>
		    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.4.1/semantic.min.css">
		    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
		    <script src="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.4.1/semantic.min.js"></script>
		</head>
		<body>
		<div style="height:60px" class="ui segment" id="version">&emsp;接口开发文档</div>
		<div class="ui text container" style="max-width:none!important;width:1200px;margin-bottom:30px;">
		    <div class="ui floating message">
		        <div class="ui grid container">
		            <div class="two wide column" style="padding:3px;background:#333">
		                <div class="ui vertical accordion menu" style="width:100%">
		                    <div class="item" style="padding:0;">
		                        <div class="content active" id="menu"></div>
		                    </div>
		                </div>
		            </div>
		            <div class="fourteen wide stretched column" style="padding:0;" id="api_list"></div>
		        </div>
		    </div>
		</div>
		<div class="ui fullscreen modal" id="api_info">
		    <div class="header">接口名称</div>
		    <div class="content">
		        <div class="ui floating message">
		            <span class='ui teal tag label' style='margin-left:15px;'></span>
		            <h4>&nbsp;请求模拟</h4>
		            <table class="ui green celled striped table">
		                <thead>
		                <tr>
		                    <th>参数</th>
		                    <th>说明</th>
		                    <th>是否必填</th>
		                    <th>其他校验规则</th>
		                    <th>默认值</th>
		                    <th>值</th>
		                </tr>
		                </thead>
		                <tbody class="params"></tbody>
		            </table>
		            <div style="display:flex;align-items:center;"></div>
		            <div class="ui fluid action input">
		                <input placeholder="请求的接口链接" type="text" class="api_url" readonly>
		                <button class="ui button blue" id="get_tk">获取TK（请先在TK文本框中输入测试手机号)</button>
		                <button class="ui button red" id="submit">请求当前接口</button>
		            </div>
		            <div class="ui blue message" id="json_output"></div>
		        </div>
		    </div>
		    <div class="actions">
		        <div class="ui cancel button">关闭</div>
		    </div>
		</div>
		
		<script>
		    var items = ` + string(bs) + `;
		    for (var module in items) {
		        if (module == "websocket") {
		            continue;
		        }
		        for (var group in items[module]) {
		            var menuName = (module == "backend" ? "管理端" : "移动端") + "." + items[module][group]["name"];
		            var menuKey = module + "." + group;
		            $("#menu").append("<a class='item' data-tab='" + menuKey + "'>" + menuName + "</a>");
		            var html_str = "<div class='ui tab' data-tab='" + menuKey + "'><table class='ui striped selectable inverted table' style='border-radius:0'><tbody>";
		            var tr_str = "";
		            $.each(items[module][group]["list"], function (key, route) {
		                tr_str += "<tr data-key='" + key + "' data-config='" + JSON.stringify(route) + "' data-method='"+route['method']+"' data-route='/" + (module == "backend" ? "backend/" : "") + (route["with_platform"] ? "{platform}/" : "") + route["route"] + "'><td>/" + (module == "backend" ? "backend/" : "") + (route["with_platform"] ? "{platform}/" : "") + route["route"] + "&nbsp;&emsp;" + route["name"] + "</td></tr>";
		            });
		            html_str += tr_str + "</tbody></table></div>";
		            $("#api_list").append(html_str);
		        }
		    }
		    $('#menu a.item').tab();
		    $('#menu a:first-child').click();
		    $('#api_list').on('click', 'tbody>tr', function () {
		        var key = $(this).data('key');
		        var route = $(this).data('route');
		        var method = $(this).data('method');
		        var config = $(this).data('config');
		        var need_auth = config.need_auth ? '&emsp;<font color="red">需要认证</font>' : ''
		        $('#api_info .header').html("<span style='font-weight:bold;'>[" + config.method.toUpperCase() + "]</span>&emsp;<span style='text-decoration:underline;'>" + route + "</span>");
		        $('#api_info .api_url').data('route',route).val(route);
		        $('#api_info .label').html(config.name + need_auth);
		        $("#json_output").html('').hide();
		        var trs = [
		            '<tr style="background:#aaa"><td>AuthToken</td><td>认证Token</td><td>' + (need_auth == '' ? '' : '<font color="red">必填</font>') + '</td><td></td><td></td><td><input id="TK" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text" placeholder="请先输入用户名，再点下方按钮获取AuthToken和RandomStr值" value="lyndon"/></td></tr>',
		            '<tr style="background:#aaa"><td>RandomStr</td><td>认证Token随机字符串</td><td>' + (need_auth == '' ? '' : '<font color="red">必填</font>') + '</td><td></td><td></td><td><input id="RS" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text" readonly/></td></tr>',
		        ];
				if(route.indexOf("{platform}") != -1) {
					trs[trs.length] = '<tr style="background:#ccc"><td>platform</td><td>平台标识</td><td><font color="red">必填</font></td><td></td><td>qp</td><td><input id="PF" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text"/></td></tr>'
				}
		        if (need_auth != "") {
		            $('#get_tk').show()
		        }
		        $.each(config.params, function (k, v) {
		            var is_must = false;
					var rules = [];
		            $.each(v["rules"], function (kk, vv) {
		                if (vv == "must") {
		                    is_must = true;
		                }else{
							rules[rules.length] = vv;
						}
		            })
		            trs[trs.length] = '<tr><td>' + v["field"] + '</td><td>' + v["label"] + '</td><td>' + (is_must ? '<font color="red">必填</font>' : '<font color="grey">非必填</font>') + '</td><td>' + rules.join(",") + '</td><td>' + (v.hasOwnProperty("default") ? v["default"] : '') + '</td><td><input name="' + v["field"] + '" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" class="C_input" type="text"/></td></tr>'
		        });
		        $('#api_info .params').html(trs.join(''))
		        $('#api_info').data({'key': key,'method': method}).modal({
		            observeChanges: true,
		            transition: 'scale',
		            centered: false,
		            closable: false,
		        }).modal('show');
		    });
		    $('#get_tk').click(function () {
		        var target = $('#TK').val();
				if(target == ""){
					$('#TK').focus();alert('用户名不能为空');return;
				}
		        $.post("/doc/get_token", {target: target}, function (data) {
		            if (data) {
		                $("#TK").val(data.token);
		                $("#RS").val(data.random_str);
		            }
		        }, "json");
		    });
			$('#api_info .params').on("keyup", "#PF", function(){
				var platform = $.trim($(this).val());
				if(platform == ""){
					platform = "{platform}";
				}
				var route = $('#api_info .api_url').data('route').replace('{platform}', platform);
				$('#api_info .api_url').val(route);
			});
		    var get_data = function () {
		        var data = {};
		        $("td .C_input").each(function (index, e) {
		            data[e.name] = $.trim(e.value);
		        });
		        return data
		    }
		    $("#json_output").hide();
		    $("#submit").on("click", function () {
		        var method = $('#api_info').data('method');
		        var key = $('#api_info').data('key');
		        var url_arr = $(".api_url").val();
		        var req_obj = {
		            url: url_arr,
		            type: method,
					cache: false,
					contentType: 'application/json',
		            beforeSend: function (XMLHttpRequest) {
		                XMLHttpRequest.setRequestHeader("Token", $('#TK').val());
		                XMLHttpRequest.setRequestHeader("RandomStr", $('#RS').val());
		            },
		            success: function (res, status, xhr) {
		                var status = xhr.status + ' ' + xhr.statusText;
		                var header = xhr.getAllResponseHeaders();
		                var json_text = JSON.stringify(res, null, 4);
		                $("#json_output").html('<pre style="white-space:pre-wrap;word-wrap:break-word;">请求返回状态 ：' + status + '<br/><hr/>请求头：<br/>' + header + '<hr/>返回数据：<br/>' + json_text + '</pre>');
		                $("#json_output").show();
		            },
		            error: function (error) {
		                $("#json_output").html('<pre style="white-space:pre-wrap;word-wrap:break-word;">请求返回状态 ：' + error.status + '<br/><hr/>错误信息：<br/>' + error.statusText);
		                $("#json_output").show();
		            }
		        }
				var data = get_data();
				if(method.toUpperCase() == "POST") {
		            req_obj.data = JSON.stringify(data);
		            req_obj.processData = false;
				}else{
					req_obj.data = data;
					req_obj.dataType = 'json';
				}
		        $.ajax(req_obj)
		    })
		</script>
		</body>
		</html>`
		template.Must(template.New("").Parse(html)).Execute(rep, nil)
	}).Methods("GET")

	Mux.HandleFunc("/doc/get_token", func(rep http.ResponseWriter, req *http.Request) {
		target := req.FormValue("target")
		token, randomStr := "", ""
		admin, err := MD("main").Select("admins", true, "id", "username='"+target+"'")
		if (err == nil) && (admin != nil) {
			id := admin.(map[string]string)["id"]
			obj := Services["auth"]
			handler := RV(obj).MethodByName("GetToken")
			if handler.IsValid() {
				tmp := handler.Call([]reflect.Value{RV("backend"), RV(id)})
				if tmp[2].IsNil() {
					token = tmp[0].Interface().(string)
					randomStr = tmp[1].Interface().(string)
				}
			}
		}
		bs, _ := json.Marshal(map[string]string{"token": token, "random_str": randomStr})
		rep.WriteHeader(200)
		rep.Header().Set("Content-Type", "application/json")
		rep.Write(bs)
	}).Methods("POST")

}
