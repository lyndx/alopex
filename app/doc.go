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
		    <title>接口文档</title>
		    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.4.1/semantic.min.css">
		    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
		    <script src="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.4.1/semantic.min.js"></script>
			<style>#api_info .ui.table td{padding:4px 10px!important}#api_info .ui.table thead th{padding:8px 10px;font-size:16px;}#api_info .ui.table,#api_info .ui.table *{border-color:#999;}</style>
		</head>
		<body>
		<div style="height:60px;font-size:20px;" class="ui segment">&nbsp;&nbsp;接口开发文档</div>
        <div class="ui grid container" style="background:#444;position:relative;width:94%!important;min-height:90%;padding-left:160px;margin:0;margin-left:30px!important;">
            <div class="two wide column" style="position:absolute;width:160px!important;height:100%;background:#333;top:0;left:0;padding:3px;background:#333">
                <div class="ui vertical accordion menu" style="width:100%">
                    <div class="item" style="padding:0;">
                        <div class="content active" id="menu"></div>
                    </div>
                </div>
            </div>
            <div class="fourteen wide stretched column" style="padding:0;width:100%!important;" id="api_list"></div>
        </div>
		<div class="ui fullscreen modal" id="api_info" style="right:0!important;">
		    <div class="header" style="font-size:20px;font-weight:bold;">接口名称</div>
		    <div class="content">
		        <div class="ui floating message">
		            <span class='ui teal tag label' style='margin-left:5px;'></span>
		            <h4>&nbsp;请求参数清单</h4>
					<table class="ui green celled striped table" style="border-radius:0;border-top:3px solid #21ba45;margin:0;">
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
					<h4>&nbsp;返回字段说明</h4>
					<table class="ui green celled striped table" style="border-radius:0;border-top:3px solid #21ba45">
		                <thead>
		                <tr>
		                    <th>参数</th>
		                    <th>说明</th>
		                </tr>
		                </thead>
		                <tbody class="return"></tbody>
		            </table>
		            <div style="display:flex;align-items:center;"></div>
					<h4>&nbsp;模拟请求</h4>
		            <div class="ui fluid action input" style="border:10px solid #ddd;">
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
		        $('#api_info .header').html("<span style='font-weight:bold;background:springgreen;padding:3px;'>[" + config.method.toUpperCase() + "]</span>&emsp;<span style='text-decoration:underline;text-underline-position:under;'>" + route + "</span>");
		        $('#api_info .api_url').data('route',route).val(route.replace('{platform}','qp'));
		        $('#api_info .label').html(config.name + need_auth);
		        $("#json_output").html('').hide();
		        var trs = [];
				if(need_auth){
					trs[trs.length] = '<tr style="background:#aaa;font-weight:bold;"><td>AuthToken</td><td>认证Token</td><td><font color="red">必填</font></td><td></td><td></td><td><input id="TK" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text" placeholder="请先输入用户名，再点下方按钮获取AuthToken和RandomStr值" value="lyndon"/></td></tr>';
		            trs[trs.length] = '<tr style="background:#aaa;font-weight:bold;"><td>RandomStr</td><td>认证Token随机字符串</td><td><font color="red">必填</font></td><td></td><td></td><td><input id="RS" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text" readonly/></td></tr>';
				}
				if(route.indexOf("{platform}") > -1) {
					trs[trs.length] = '<tr style="background:#999;font-weight:bold;"><td>platform</td><td>平台标识</td><td><font color="red">必填</font></td><td></td><td>qp</td><td><input id="PF" value="qp" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text"/></td></tr>';
				}
                if(route.indexOf("common/upload/{path}") > -1) {
					trs[trs.length] = '<tr style="background:#999;font-weight:bold;"><td>path</td><td>上传目录</td><td><font color="red">必填</font></td><td></td><td></td><td><input id="PATH" value="" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;" type="text"/></td></tr>';
					trs[trs.length] = '<tr style="font-weight:bold;"><td>file</td><td>文件</td><td><font color="red">必填</font></td><td></td><td></td><td><input class="C_input" type="file" name="file" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;"/></td></tr>';
				}
		        if (need_auth != "") {
		            $('#get_tk').show();
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
		            });
		            trs[trs.length] = '<tr style="font-weight:bold;">'+
										'<td>' + v["field"] + '</td>'+
										'<td>' + v["label"] + '</td>'+
										'<td>' + (is_must ? '<font color="red">必填</font>' : '<font color="grey">非必填</font>') + '</td>'+
										'<td>' + rules.join(",") + '</td>'+
										'<td>' + (v.hasOwnProperty("default") ? v["default"] : '') + '</td>'+
										'<td><input class="C_input" type="text" name="' + v["field"] + '" style="width:100%;padding:5px 10px;border:1px solid #ccc;outline:0;border-radius:4px;"/></td>'+
									  '</tr>';
		        });
		        $('#api_info .params').html(trs.join(''))
				var rts = [];
				for(var field in config.return) {
		            rts[rts.length] = '<tr style="font-weight:bold;"><td>' + field + '</td><td>' + config.return[field] + '</td></tr>'
		        }
		        $('#api_info .return').html(rts.join(''))
		        $('#api_info').data({'key': key,'method': method}).modal({
		            observeChanges: true,
		            transition: 'scale',
		            centered: true
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
			$('#api_info .api_url').val(($('#api_info .api_url').data('route')+'').replace('{platform}', 'qp'));
			var change_route = function(){
				var platform = $.trim($('#PF').val());
				if(platform == ''){
					platform = '{platform}';
				}
				var path = $.trim($('#PATH').val());
				if(path == ''){
					path = '{path}';
				}
				var route = ($('#api_info .api_url').data('route')+'').replace('{platform}', platform).replace('{path}', path);
				$('#api_info .api_url').val(route);
			};
			$('#api_info .params').on('keyup', '#PF', change_route).on('keyup', '#PATH', change_route);
		    var get_data = function () {
				if ($('#api_info .api_url').data('route').indexOf("common/upload/{path}") > -1) {
					let data = new FormData();
					$("td .C_input").each(function (i, e) {
			            data.append(e.name, (e.name == 'file') ? (e.files.length > 0 ? e.files[0] : '') : $.trim(e.value));
			        });
					return data;
		        }
				var data = {};
		        $("td .C_input").each(function (index, e) {
		            data[e.name] = $.trim(e.value);
		        });
		        return data
		    };
		    $("#json_output").hide();
		    $("#submit").on("click", function () {
		        var method = $('#api_info').data('method');
		        var key = $('#api_info').data('key');
		        var url_arr = $(".api_url").val();
		        var req_obj = {
		            url: url_arr,
		            type: method,
					cache: false,
					dataType:'json',
		            beforeSend: function (XMLHttpRequest) {
		                XMLHttpRequest.setRequestHeader("Debug", 'true');
		                XMLHttpRequest.setRequestHeader("Token", $('#TK').val());
		                XMLHttpRequest.setRequestHeader("Random_str", $('#RS').val());
		            },
		            success: function (res, status, xhr) {
		                var data_text = JSON.stringify(res.data, null, 4);
		                var params_all_text = JSON.stringify(res.request.params, null, 4);
		                var params_needed_text = JSON.stringify(res.request.needed, null, 4);
		                $("#json_output").html('<pre style="white-space:pre-wrap;word-wrap:break-word;font-weight:bold;">请求返回状态 ：' + res.code + '<br/><hr/>请求返回消息 ：' + res.message + (res.message_detail ? ('('+res.message_detail+')') : '') + '<br/><hr/>有效请求参数：<br/>' + params_needed_text + '<br/><hr/>所有请求参数：<br/>' + params_all_text + '<br/><hr/>返回数据：<br/>' + (data_text ? data_text : 'Empty . . . .') + '</pre>');
		                $("#json_output").show();
		            },
		            error: function (error) {
		                $("#json_output").html('<pre style="white-space:pre-wrap;word-wrap:break-word;font-weight:bold;">请求返回状态 ：' + error.status + '<br/><hr/>返回错误信息：<br/>' + error.statusText);
		                $("#json_output").show();
		            }
		        }
				if ($('#api_info .api_url').data('route').indexOf("common/upload/{path}") > -1) {
					req_obj.contentType = false;
				}
				var data = get_data();
				if(method.toUpperCase() == "POST") {
					req_obj.data = data;
					if ($('#api_info .api_url').data('route').indexOf("common/upload/{path}") == -1) {
			            req_obj.data = JSON.stringify(data);
					}
		            req_obj.processData = false;
				}else{
					req_obj.data = data;
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
