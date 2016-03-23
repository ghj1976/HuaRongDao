
/*
  特色: 经典布局|数字字母|八卦|异类|没过关记录的布局|我没玩过的布局
  难度：入门(<30步)| 进阶(<50步)| 中阶(<80步)| 高阶(<100步)| 超阶(<120步)| 无敌(>=120步)
  编号: 1-50 | 51-100 | 101-150 | 151-200 | ... | 所有布局(无图)
  横数: 0横|1横|2横|3横|4横|5横
  曹操位置: 顶线正中|顶线|二线|三线|底线

  "s_" 开头表示server端初始化的数据
  全局变量:s_gates = [
    {
      //from server
      "id":gate_id,
      "name":gate_name,
      "rmin":min_step,
      "layout":layout,
      "rcnt":played_by_users,

      //script generated
      "iplayed":played_by_current_user,
      "bing_count":bing_count,//count of 小兵
      "hbar_count":hbar_count,//count of 横将
      "cao_y":caocao_seat, //from top to bottom:1,2,3,4
    },
  ]
  全局变量:s_iplayed = [id,id,id,...]
  全局变量:g_id2gate = {1:obj_gate_in_s_gates,...}
*/

//!!!!change the global variable
function generate_extra_gates_attrs() {
  g_id2gate = {};
  var a=3;
  var i=0;
  var gate = null;
  for (i=0; i<s_gates.length; i++) {
    gate = s_gates[i];
    g_id2gate[gate.id] = gate;
    gate["layout_str"] = gate.layout;
    gate.layout = eval(gate.layout_str); //from str to array

    var layout = gate.layout;
    var cao_y = 0;
    var bing_count = 0;
    var hbar_count = 0;
    for (var j=0; j<layout.length; j+=3) {
      if (layout[j] == 1) bing_count += 1;
      else if (layout[j] == 3) hbar_count += 1;
      else if (layout[j] == 4) cao_y = layout[j+2];
    }
    gate["cao_y"] = cao_y;
    gate["bing_count"] = bing_count;
    gate["hbar_count"] = hbar_count;
  }
  for (i=0; i<s_iplayed.length; i++) {
    gate = g_id2gate[s_iplayed[i]];
    gate["iplayed"] = true;
  }
}

function group_special() {
  //特色: 经典布局|数字字母|八卦|异类|没过关记录的布局|我没玩过的布局
  var group = [
    {"name":"经典布局",
     "gates":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,39,40,41,42,46,47,48,49,50,51,52,53,55]},
    {"name":"数字字母",
     "gates":[84,85,86,87,88,89,90,91,92,93]},
    {"name":"伏羲八卦",
     "gates":[100,104,129,130,131,132,133,134,135,136,137,138,139,140,141,142,143,144,145,146,147,148,149,150,151,152,153,154,155,156,157,158,159,160,161,162,163,164,165,166,167,168,169,170,171,172,173,174,175,176,177,178]},
    {"name":"定式",
     "gates":[187,188,189,190,191,192,193,194,195,196,197,198,199]}
  ];

  var inot_played = [];
  var i_played = [];
  var noone_played = [];
  var abnormals = [];
  var played = [];
  var i = 0;
  for (i=0; i<s_gates.length; i++) {
    var gate = s_gates[i];
    played.push({"id":gate.id, "count":gate.rcnt});
    if (!gate.iplayed) inot_played.push(gate.id);
    else i_played.push(gate.id);
    if (!gate.rcnt) noone_played.push(gate.id);
    if (gate.bing_count != 4) abnormals.push(gate.id);
  }
  played.sort(function(a,b) {return a.count - b.count;});
  var rare_played = [];
  for (i=0; i<Math.min(30, played.length); i++) {
    rare_played.push(played[i].id);
  }
  group.push({"name":"我没玩过的", "gates":inot_played});
  group.push({"name":"我玩过的", "gates":i_played});
  group.push({"name":"很少人玩的", "gates":rare_played});
  group.push({"name":"没人过关的", "gates":noone_played});
  group.push({"name":"非标准布局", "gates":abnormals});
  return group;
}

function group_hard_levels() {
  /* var names = [ "简单(<20步)", "入门(<40步)", "进阶(<60步)", "中阶(<80步)",
                "高阶(<100步)", "超阶(<120步)", "无敌(>=120步)", "未知" ];
  */
  var names = [ "简单", "入门", "进阶", "中阶", "高阶", "超阶", "无敌", "未知"];
  var group = {};
  var i = 0;
  var name = "";
  for (i=0; i<s_gates.length; i++) {
    var gate = s_gates[i];
    var bs = gate.rmin;
    name = names[Math.floor(bs/20)];
    if (bs === 0) name = names[names.length - 1];
    if (group[name]) group[name].push(gate.id);
    else group[name] = [gate.id];
  }
  var ga = [];
  for (i=0;i<names.length;i++) {
    name = names[i];
    ga.push({"name":name, "gates":group[name]});
  }
  return ga;
}

function group_gate_ids() {
  var group = {};
  var names = [];
  var i = 0;
  var name = "";
  for (i=0; i<s_gates.length; i+=50) {
    var from = (Math.floor(i/50)*50+1);
    var to = (Math.floor(i/50) + 1)*50;
    if (to > s_gates.length)
      to = s_gates.length;
    name = from + '-' + to;
    names.push(name);
    group[name] = [from];
    for (var j=from; j<to; j++) {
      group[name].push(s_gates[j].id);
    }
  }
  var ga = [];
  for (i=0;i<names.length;i++) {
    name = names[i];
    ga.push({"name":name, "gates":group[name]});
  }
  return ga;
}

function group_hbars() {
  var group = {};
  var names = ['0横', '1横', '2横', '3横', '4横', '5横'];
  var i = 0;
  var name = "";
  for (i=0; i<s_gates.length; i++) {
    var gate = s_gates[i];
    name = names[gate.hbar_count];
    if (group[name]) group[name].push(gate.id);
    else group[name] = [gate.id];
  }
  names.sort();
  var ga = [];
  for (i=0;i<names.length;i++) {
    name = names[i];
    ga.push({"name":name, "gates":group[name]});
  }
  return ga;
}

function group_cao_seats() {
  var group = {};
  var names = ['1线', '2线', '3线', '4线'];
  var i = 0;
  var name = "";
  for (i=0; i<s_gates.length; i++) {
    var gate = s_gates[i];
    name = names[gate.cao_y];
    if (group[name]) group[name].push(gate.id);
    else group[name] = [gate.id];
  }
  names.sort();
  var ga = [];
  for (i=0;i<names.length;i++) {
    name = names[i];
    ga.push({"name":name, "gates":group[name]});
  }
  return ga;
}

function set_up_group() {
  generate_extra_gates_attrs();
  var hrd_group = {};
  hrd_group["特殊分类"] = group_special();
  hrd_group["难度分类"] = group_hard_levels();
  hrd_group["布局编号"] = group_gate_ids();
  hrd_group["横数分类"] = group_hbars();
  hrd_group["曹操位置"] = group_cao_seats();
  return hrd_group;
}

//global variable!!!
//  g_groups = [[xx,xx], [xxx,xx,xxx],...] array of id-arrays
//
function generate_group_filters() {
  var groups = set_up_group();
  g_groups = []; //!!!
  var html = "";
  for (var k in groups) {
    html += "<li><strong>" + k + "</strong>: ";
    var group = groups[k];
    var a=0;
    for (var i=0; i<group.length; i++) {
      var s = group[i];
      var gate_count = 0;
      if (s.gates) gate_count = s.gates.length;
      html += '<a href="#hrd-filters" title="共有' + gate_count + '个开局" ';
      html += 'onclick="group_filter(' + g_groups.length + ', this);return false">' + s.name + '</a> | ';
      g_groups.push(s.gates);
    }
    html = html.substring(0, html.length - 2);
    html += "</li>";
  }
  $("#hrd-filters").html(html);
}

function group_filter(group_id, obj) {
  var group = g_groups[group_id];
  if (!group) return;
  var html = "<tr>"; //浪费第一个<tr></tr>就浪费吧，逻辑简单些
  for (var i=0; i<group.length; i++) {
    var gate_id = group[i];
    var gate = g_id2gate[gate_id];
    if (!gate) continue; //防止布局编号出错带来问题
    if (i % 5 === 0) html += "</tr><tr>";
    html += '<td valign="top"><center><div class="layoutcanvas" tag="' + gate.layout_str;
    html += '" id="layout' + gate_id + '"';
    html += '" onclick=\'window.location="/youxi/hrd/' + gate_id +'/";return false\'>';
    html += '<div class="layoutbottombar"> </div></div></center><br/>#' + gate_id;
    html += ':<span style="font-weight:bold"><a href="/youxi/hrd/' + gate_id + '/#spec">';
    html += gate.name + '</a></span>';
    if (gate.rcnt) {
      html += '<br/><small><a href="/youxi/hrd/' + gate_id + '/replays/#gate-spec">过关记录';
      html += gate.rcnt+ '个</a>,最少' + gate.rmin + '步</small>';
    }
    html += '<br/><small>by <a href="/youxi/user/' + gate.uid + '/">';
    html += gate.un + '</a></td></small>';
  }
  html += "</tr>";
  $("#hrd-gates").html(html);
  //render on canvas
  hrd_render_on_canvas();
  $("#hrd-filters a").removeClass("selected");
  if (obj)
    $(obj).addClass("selected");
  else
  {
    $("#hrd-filters a:first").addClass("selected");
  }
}

