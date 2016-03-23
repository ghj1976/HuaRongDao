//----------------------------
//     0   1   2   3
//    --- --- --- ---
// 0 | 2a  4   4   2c
// 1 | 2a  4   4   2c
// 2 | 2b  3   3   2d
// 3 | 2b  1b  1c  2d
// 4 | 1a  0   0   1d
//----------------------------
function Cell(type, x, y) { this.x = x; this.y = y; this.type = type; }

/* return true if cells's status is in "WIN" state */
function has_win() {
  for (var i=0; i<cells.length; i++) {
    var cell = cells[i];
    if (cell.type == 4)
      return cell.x == 1 && cell.y == 3;
  }
  return false;
}

function save_gate_result() {
  var myurl = "/youxi/hrd/" + gate_id + "/win/";
  //结果里面只要存block的index以及它的移动方向的index就行
  //要求: server和client端的directions设计要保持一致
  var params = [];
  for (var i=0; i<game_steps.length; i++) {
    var step = game_steps[i];
    var index = cellid_to_index[step[0]]; //game_steps.push([obj.id, move]);
    params.push(index);
    params.push(step[1]);
  }
  param_str = "steps=[" + params.join(",") + "]";

  $.post(myurl, param_str);
}

/* show the "WIN" result to user */
function show_win_and_change_status() {
  save_gate_result();

  locked = true; //lock until been reset
  soundManager.play("win");

  var str= "恭喜过关! 共" + get_step_count() + "步，";
  var step_count = get_step_count();
  if (step_count <= 100)
    str += "大侠~太牛了!";
  else if (step_count <= 200)
    str += "厉害啊~!";
  else if (step_count <= 500)
    str += "不错，加油!";
  else
    str += "下次快些哦!";
  show_info("刷新 或 点击“重新开始”重玩本局");
  alert(str);
}

/* return true if cell has conflict with existing matrix */
function have_conflict(matrix, cell) {
  var tests = [
      [],
      [[0, 0]],
      [[0, 0], [0, 1]],
      [[0, 0], [1, 0]],
      [[0, 0], [1, 0], [0, 1], [1, 1]]];
  var test = tests[cell.type];
  for (var i=0; i<test.length; i++) {
    var t = test[i];
    var x = cell.x + t[0];
    var y = cell.y + t[1];
    if (x < 0 || x > 3) {
      return true;
    }
    if (y < 0 || y > 4) {
      return true;
    }
    if (matrix[x][y]) {
      return true;
    }
    matrix[x][y] = cell.type;
  }
  return false;
}

/* return true if target obj can move on the given direction(dir_index) */
function can_move(obj, dir_index) {
  /* judge if we can move on current dir_index */
  /* try to fill the zero-table, if we hit an "non-zero" then it's forbidden */
  var dx = directions[dir_index][0];
  var dy = directions[dir_index][1];
  var target_cell = cell_div_map[obj.id];
  var x = target_cell.x;
  var y = target_cell.y;
  var test_cell = new Cell(target_cell.type, x + dx, y + dy);
  var matrix = [[0,0,0,0,0], [0,0,0,0,0], [0,0,0,0,0], [0,0,0,0,0]];
  //catch out of bounds issue
  if (have_conflict(matrix, test_cell)) return false;
  for (var i=0; i<cells.length; i++) {
    var cell = cells[i];
    if (cell != target_cell && have_conflict(matrix, cell))
      return false;
  }
  return true;
}

function show_info(str) {
  $("#infobox").html(str);
}

function update_recent_steps() {
  var length = game_steps.length;
  if (length <= 0) {
    show_info("---- life is short, play more ----");
  } else {
    var step = game_steps[length - 1];
    var str = $("#" + step[0]).attr("tag") + directions[step[1]][2];
    show_info(str);
  }
  $("#countbox").text(get_step_count());
}

function mouse_down(e, obj) {
  if (locked) {
    if (typeof(is_replay_mode) != "undefined" && is_replay_mode)
      show_info('这是结果回放页面，想玩点击<a href="/youxi/hrd/' + gate_id + '/#gate-spec">这里</a>');
    else if (typeof(is_in_add_page) == "undefined")
      alert("点击“重新开始”或刷新此页面以重玩本布局");
    return; //do nothing if the block is locked
  }

  var count = 0;
  var move = null;
  var dir_next = 0;
  if (game_steps.length > 0 && game_steps[game_steps.length - 1][0] === obj.id) {
    dir_next = (dir_revert[last_dir_index] + 1) % 4;
  }

  while (count < 4) { //4 directions
    if (can_move(obj, dir_next)) {
      move = dir_next;
      break;
    }
    dir_next = (dir_next + 1) % 4;
    count++;
  }

  if (move === null ) {
    //no valid move for current block
    soundManager.play("nomove");
    return;
  }
  last_dir_index = move;

  //let's move!
  //  we could get at MOST 2 moves!
  //move the cell
  var left = parseInt($(obj).css('left'), 10);
  var top = parseInt($(obj).css('top'), 10);
  var direction = directions[move];
  var moves = ["-=", "", "+="];
  var premove = moves[direction[0] + 1];
  if (premove)
    $(obj).animate({left: premove+unit }, 100);
  premove = moves[direction[1] + 1];
  if (premove)
    $(obj).animate({top: premove+unit }, 100);

  soundManager.play("move");

  //save the cell
  cell_div_map[obj.id].x += direction[0];
  cell_div_map[obj.id].y += direction[1];

  //save last move
  game_steps.push([obj.id, move]);
  update_recent_steps();

  //have we win?
  if (has_win()) {
    show_win_and_change_status();
  }
}

/*当前的设计可能会在扩展华容道配置(不是十个方块)的时候存在潜在问题，需要修改
 */
function create_cell_divs() {
  var div_arr = ['11', '12', '13', '14', '15', '16', '17', '18',
                 '19', '110', '111', '112', '113', '114', '115', '116',
                 '21', '22', '23', '24', '25', '26',
                 '31', '32', '33', '34', '35', '36',
                 '41'];
  var name_arr = ['兵', '卒', '勇', '士', '兵2', '卒2', '勇2', '士2',
                  '兵3', '卒3', '勇3', '士3', '兵4', '卒4', '勇4', '士4',
                  '赵云', '黄忠', '马超', '张飞', '关羽', '魏延',
                  '关羽', '张飞', '马超', '黄忠', '赵云', '魏延',
                  '曹操'];
  var jqobj = $("#canvas");
  for (var i=0; i<div_arr.length; i++) {
    var temp = div_arr[i];
    var html_str = '<div id="hrdcell' + temp + '" class="hrdcell"';
    html_str += ' onmousedown="mouse_down(event, this)"';
    html_str += ' tag="' + name_arr[i] + '"> </div>';
    jqobj.append(html_str);
  }
}

/* this function CHANGES cell_div_map & cellid_to_index */
function hrd_render(cells) {
  var ccounts = [0, 0, 0, 0, 0];
  for (var i=0; i<cells.length; i++) {
    var cell = cells[i];
    var x = cell.x * unit + borderX;
    var y = cell.y * unit + borderY;
    ccounts[cell.type] += 1;
    var id_index = ccounts[cell.type];
    if ((cell.type == 2 || cell.type == 3) && ccounts[2] + ccounts[3] > 5)
      id_index = 6;
    var id ="hrdcell" + cell.type + id_index;
    var jqobj = $("#" + id);
    jqobj.show().css('left', x + 'px').css('top', y + 'px');
    cell_div_map[id] = cell;
    cellid_to_index[id] = i;
    index_to_cellid[i] = id;
  }
}

/* ab - Action Buttion functions */

function ab_rollback() {
  if (game_steps.length <= 0) {
    soundManager.play("nomove");
    return;
  }

  soundManager.play("move");
  var step = game_steps.pop();
  var id = step[0];
  var move = step[1];
  var direction = directions[move];
  //roback the block
  var jqobj = $("#" + id);
  var left = parseInt(jqobj.css('left'), 10);
  var top = parseInt(jqobj.css('top'), 10);
  jqobj.css('left', left - unit * direction[0]);
  jqobj.css('top', top - unit * direction[1]);
  //roback the cell map
  cell_div_map[id].x -= direction[0];
  cell_div_map[id].y -= direction[1];

  update_recent_steps();
}

function ab_reset() {
  if (game_steps.length <= 0) {
    soundManager.play("nomove");
    return;
  }
  if (!locked) { //lock状态下就无需提示了
    if (!confirm("您已经走了几步了，真的要重新开始吗？")) return;
  }

  locked = false;

  soundManager.play("move");
  //reset all the status
  last_dir_index = 0;
  var cstr = $("#canvas").attr('tag');
  carr = eval(cstr);
  for (var i=0; i<carr.length; i+=3) { cells[i / 3] = new Cell(carr[i], carr[i+1], carr[i+2]); }
  game_steps = [];
  hrd_render(cells);

  update_recent_steps();
}

function State() {
  this.last_dir_index = last_dir_index;
  this.cells = cells.clone();
  this.step_count = get_step_count();
  this.game_steps = game_steps.clone();
}
function ab_save() {
  last_state = new State();
  soundManager.play("save");
}

function ab_restore() {
  if (!confirm("真的要恢复到第 " + last_state.step_count + " 步吗？")) return;

  last_dir_index = last_state.last_dir_index;
  cells = last_state.cells.clone();
  game_steps = last_state.game_steps.clone();

  hrd_render(cells);
  update_recent_steps();

  soundManager.play("restore");
}

function get_step_count() {
  var count = 0;
  var last_id = null;
  for (var i=0; i<game_steps.length; i++) {
    var id = game_steps[i][0];
    if (id == last_id) continue;
    last_id = id;
    count++;
  }
  return count;
}

function initializations() {
  if (typeof(page_initilized) != "undefined" && page_initilized) return;
  page_initilized = true;

  /*globals for HRD game*/
  //game canvas border
  borderX = 20;
  borderY = 60;
  unit = 80; //min-cell size
  /*moving directions for cell*/
  directions = [[0, 1, "向下"], [0, -1, "向上"], [1, 0, "向右"], [-1, 0, "向左"]];
  dir_revert = [1, 0, 3, 2]; //相反的方向
  last_dir_index = 0; //0, 1, 2, 3
  /*html div map to cell*/
  cell_div_map = {};
  /*html div id to cell index*/
  cellid_to_index = {};
  /*cell index to html div id*/
  index_to_cellid = {};
  /*all cells*/
  cells = [];
  /*game steps*/
  game_steps = [];

  create_cell_divs();

  //load cells
  /*var cstr = '[2,0,0, 4,1,0, 2,3,0, 2,0,2, 3,1,2, 2,3,2, 1,0,4, 1,1,3, 1,2,3, 1,3,4]';*/
  //
  //----------------------------
  //     0   1   2   3
  //    --- --- --- ---
  // 0 | 2a  3   3   2c
  // 1 | 2a  1b  1c  2c
  // 2 | 2b  4   4   2d
  // 3 | 2b  4   4   2d
  // 4 | 1a  0   0   1d
  //----------------------------
  /*var cstr = '[2,0,0, 3,1,0, 2,3,0, 1,1,1, 1,2,1, 2,0,2, 4,1,2, 2,3,2, 1,0,4, 1,3,4]';*/
  var cstr = $("#canvas").attr('tag');
  carr = eval(cstr);

  //set the sound manager swf file position
  //NOTE the trailing slash is REQUIRED!
  soundManager.url = '/static/flash/';
  soundManager.debugMode = false;
  soundManager.flashVersion = 9;

  var sounds = [ "hover", "move", "nomove",
                 "reset", "restore", "save", "msg", "win" ];

  soundManager.onload = function() {
    for (var i=0; i<sounds.length; i++) {
      soundManager.createSound(sounds[i], '/static/sound/youxi/' + sounds[i] + '.mp3');
    }
    soundManager.play("msg");

    //set up sound for "ab" buttons
    $("#operations li").hover(function(e){soundManager.play("hover");},function(){});
  };
}

function hrd_play_creation() {
  initializations();
  cells = [];
  for (var i=0; i<carr.length; i+=3) { cells[i / 3] = new Cell(carr[i], carr[i+1], carr[i+2]); }
  hrd_render(cells);
  //save init state
  init_state = new State();
  last_state = new State();
}

/* hrd_home stuff */
Cell.prototype.clone = function() { return new Cell(this.type, this.x, this.y); };
function get_cell_w(cell) { return (cell.type > 2) ? 2 : 1; }
function get_cell_h(cell) { return (cell.type % 2) ? 1 : 2; }
function get_cell_color(cell) {
  return ["rgb(0,200,0)", "rgb(0,0,200)", "rgb(200,0,200)", "rgb(200,0,0)"][cell.type - 1];
}

/* hrd_small_canvas render and toggle for list/add/ */
function clear_marks() {
  set_cookie("hrdlist_ids", "");
  $(".layout").removeClass("marked");
}
function toggle_mark(layout) {
  $(layout).toggleClass("marked");
}
function toggle_mark_and_remember(layout) {
  var ids = get_cookie("hrdlist_ids");
  if (!ids)
    ids = "";
  var id = layout.id.substring("layout".length); //layout{{#id}}
  var pos = ("," + ids + ",").indexOf("," + id + ",");
  if (pos == -1) {
    ids += "," + id;
  } else {
    ids = ids.substring(0, pos) + ids.substring(pos + ("," + id).length);
  }
  set_cookie("hrdlist_ids", ids);
  toggle_mark(layout);
}

function hrd_render_on_canvas() {
  var unit = 20;
  var border = 1;
  var i = 0;
  var ids = get_cookie("hrdlist_ids");
  $('.layoutcanvas').each(function() {
    //calc the cells
    var layout = eval($(this).attr('tag'));
    var cells = [];
    for (i=0; i<layout.length; i+=3) {
      cells.push(new Cell(layout[i], layout[i+1], layout[i+2]));
    }
    //setting marks
    this.oncontextmenu = function() { return false; };
    $(this).mousedown(function(eve) { if (eve.button & 2) toggle_mark_and_remember(this); });
    //toggle mark for list/add/ page
    var id = this.id.substring("layout".length); //layout{{#id}}
    var pos = ("," + ids + ",").indexOf("," + id + ",");
    if (pos != -1)
      toggle_mark(this);
    //actually render it
    var canvas = $(this);
    for (i=0; i<cells.length; i++) {
      var cell = cells[i];
      var x = cell.x * unit + border;
      var y = cell.y * unit + border;
      var newdiv = $('<div class="celltype' + cell.type + '"> </div>');
      newdiv.css('left', x);
      newdiv.css('top', y);
      canvas.append(newdiv);
    }
  });
}

////////////////////////////////////////////////////////////////////////////////////////////////
// ------ hrd_replay stuff ------
//
function set_replay(obj, rid_or_steps, step_count) {
  hrd_replay_home();

  var steps = null;
  if (obj !== null) {
    steps = eval($(obj).attr("tag"));
    $(".replay").removeClass("selected-replay");
    $(obj).addClass("selected-replay");
    $("#replay-info").text("#" + rid_or_steps + "," + step_count + "步");
  } else {
    $("#replay-info").text(step_count + "步");
    steps = rid_or_steps;
  }
  hrd_replay_steps = [];
  var i = 0;
  for (i=0; i<steps.length; i+=2) {
    hrd_replay_steps[i / 2] = [steps[i], steps[i+1]];
  }
  hrd_replay_index = 0;

  //calculate the last step status and the step-index-mapping-to-step-count
  hrd_replay_step2count = {0:0};
  hrd_replay_laststatus = [];
  for (i=0; i<cells.length; i++) {
    hrd_replay_laststatus.push([cells[i].x, cells[i].y]);
  }

  var count = 0;
  var last_index = null;
  for (i=0; i<hrd_replay_steps.length; i++) {
    var index = hrd_replay_steps[i][0];
    if (index != last_index) {
      last_index = index;
      count++;
    }
    hrd_replay_step2count[i + 1] = count;
    var direction = directions[hrd_replay_steps[i][1]];
    hrd_replay_laststatus[index][0] += direction[0];
    hrd_replay_laststatus[index][1] += direction[1];
  }
}

/* ------ replay control funcions ------ */
function hrd_replay_movenext() {
  if (hrd_replay_index >= hrd_replay_steps.length) {
    hrd_replay_pause();
    show_info("已经是最后一步了");
    return false;
  }

  var step = hrd_replay_steps[hrd_replay_index];
  var jqobj = $("#" + index_to_cellid[step[0]]);
  //move the obj
  var left = parseInt(jqobj.css('left'), 10);
  var top = parseInt(jqobj.css('top'), 10);
  var direction = directions[step[1]];
  jqobj.css('left', left + unit * direction[0]);
  jqobj.css('top', top + unit * direction[1]);

  hrd_replay_index++;

  $("#countbox").text(hrd_replay_step2count[hrd_replay_index]);

  soundManager.play("move");
  return true;
}

function hrd_replay_rollback() {
  if (hrd_replay_index === 0) return false;

  hrd_replay_index--;

  var step = hrd_replay_steps[hrd_replay_index];
  var jqobj = $("#" + index_to_cellid[step[0]]);
  //move the obj back
  var left = parseInt(jqobj.css('left'), 10);
  var top = parseInt(jqobj.css('top'), 10);
  var direction = directions[step[1]];
  jqobj.css('left', left - unit * direction[0]);
  jqobj.css('top', top - unit * direction[1]);

  $("#countbox").text(hrd_replay_step2count[hrd_replay_index]);

  soundManager.play("move");
  return true;
}

function hrd_replay_play() {
  if (typeof(replay_interval) != "undefined" && replay_interval !== null) {
    clearInterval(replay_interval);
    replay_interval = null;
  }
  replay_interval = setInterval('hrd_replay_movenext()', 500);
  show_info("开始播放，一秒两步");
}

function hrd_replay_pause() {
  if (typeof(replay_interval) != "undefined" && replay_interval !== null) {
    clearInterval(replay_interval);
    replay_interval = null;
    show_info("播放过程已暂停");
  }
}

function hrd_replay_prev() {
  hrd_replay_pause();
  hrd_replay_rollback();
}

function hrd_replay_next() {
  hrd_replay_pause();
  hrd_replay_movenext();
}

function hrd_replay_home() {
  hrd_replay_pause();

  hrd_replay_index = 0;

  //reset the status
  for (var i=0; i<cells.length; i++) {
    var jqobj = $("#" + index_to_cellid[i]);
    var cell = cells[i];
    var x = cell.x * unit + borderX;
    var y = cell.y * unit + borderY;
    jqobj.show().css('left', x + 'px').css('top', y + 'px');
  }

  if (typeof(hrd_replay_step2count) != "undefined") {
    $("#countbox").text(hrd_replay_step2count[hrd_replay_index]);
    soundManager.play("move");
  }
}

function hrd_replay_end() {
  hrd_replay_index = hrd_replay_steps.length;

  //set the status to end
  for (var i=0; i<hrd_replay_laststatus.length; i++) {
    var jqobj = $("#" + index_to_cellid[i]);
    var dxdy = hrd_replay_laststatus[i];
    var x = dxdy[0] * unit + borderX;
    var y = dxdy[1] * unit + borderY;
    jqobj.show().css('left', x + 'px').css('top', y + 'px');
  }

  $("#countbox").text(hrd_replay_step2count[hrd_replay_index]);

  soundManager.play("move");
}

function layout_changed() {
    hrd_replay_home();
    $("#step").text("0");
    $("#replay-info").text("");
    $("#calcgate").hide();
    $("#controlbox").hide();
    show_info("完成布局以后请进行验证");
    $("#tick").hide();
}
function init_add_page() {
  locked = false;
  is_in_add_page = true;
  $("#updatelayout").click(function() { update_layout(); });
  $("#calcgate").click(function() { calc_gate(); });
  $("#layoutbox").keypress(layout_changed);
  $("#layoutbox").change(layout_changed);

  hrd_play_creation();
}

function update_layout() {
  //update the layout
  $("#controlbox").hide();
  $("#step").text("");
  $("#countbox").text("");
  $("#replay-info").text("");
  var ss = $("#layoutbox").val().toUpperCase().split("\n");
  var b = [0, 0, 0, 0];
  var c = 0;
  var d = 0;
  var result = [];
  for (var i=0; i<Math.min(ss.length, 5); i++) {
    var s = ss[i];
    s = s.replace(/[^DVHB ]/g, "");
    for (var j=0; j<Math.min(s.length, 4); j++) {
      var n = s.charAt(j);
      if (n == 'D')
        result.push(1, j, i);
      else if (n == 'V') {
        if (i === 0 || b[j] === 0) {
          b[j] = 2;
          result.push(2, j, i);
        } else
          b[j] = 0;
      } else if (n == 'H') {
        if (c === 0) {
          result.push(3, j, i);
          c = 3;
        } else
          c = 0;
      } else if (n == 'B') {
        if (d % 4 === 0) {
          result.push(4, j, i);
          d = 1;
        } else
          d += 1;
      }
    }
  }
  if (!is_valid_layout(result)) {
    alert("排列有误，请修改以后再提交验证");
    locked = true;
    return;
  }

  hrd_hide_cells();
  $("#tick").show();
  show_info("布局通过验证，可进行最短路径计算");
  if (!document.getElementById("add_rollback"))
    $('<a id="add_rollback" style="position:absolute" href="#" onclick="ab_rollback();return false">回退一步</a>').insertBefore("#controlbox");
  locked = false;
  game_steps = [];

  carr = result;
  hrd_play_creation();

  $("#calcgate").show();
}

function is_valid_layout(result) {
  //只要有曹操在中间，而且两个空的就行了，如果全是小块，就是14*1=15个
  return (result.length >= 24 && result.length <= 45);
}

function hrd_hide_cells() {
  for (var i=0; i<cells.length; i++) {
    var id = index_to_cellid[i];
    $("#" + id).hide();
  }
  cell_div_map = {};
  cellid_to_index = {};
  index_to_cellid = {};
}

function calc_gate() {
  update_layout();
  var layout_str = carr.join("");
  layout_str += $("#winrule").val();
  $("#calcgate").attr("disabled", "true");
  $.ajax({
    type: "POST",
    dataType: 'text',
    url: "/youxi/hrd/add/",
    data: {'layout':layout_str},
    success: function(msg){
      $("#calcgate").removeAttr("disabled");

      if (msg.length <= 0) {
        alert("自动求解失败，请仔细检查布局是否合理或联系管理员");
        return;
      }
      if (msg.charAt(0) != "[") {
        alert(msg);
        return;
      }
      if (document.getElementById("add_rollback"))
        $("#add_rollback").remove();

      $("#controlbox").show();
      var arr = eval(msg);
      var step_count = arr[0];
      set_replay(null, arr.slice(1), step_count);
      locked = true;
      alert("自动求解成功，该布局最快需要" + step_count + "步!\n如果您觉得该布局不错，欢迎联系管理员提交布局!");
    },
    error: function(msg){
      //@@@@
      $("#calcgate").removeAttr("disabled");
    }
  });
}

$(document).ready(function() {
  var gatenumobj = document.getElementById("select_gate_number");
  if (gatenumobj) { //快速选关
    $(gatenumobj).keydown(function(eve){
      var target_id = parseInt($(gatenumobj).val(), 10);
      if (eve.target != this) return;
      if (eve.keyCode == 13)
        window.location.href = 'http://fayaa.com/youxi/hrd/' + target_id + "/";
    });
    $("#select_gate_button").click(function(eve){
      var target_id = parseInt($(gatenumobj).val(), 10);
      window.location.href = 'http://fayaa.com/youxi/hrd/' + target_id + "/";
    });
    $.get("/youxi/hrd/gates/", function(data) {
      eval(data); //here we got 'gates'
      var gatestr = "";
      for (var i=0; i<gates.length; i++) {
        var gate = gates[i];
        gatestr += '<option value="' + gate.id + '"';
        if (gate_id.toString() == gate.id.toString()) {
          gatestr += ' selected="selected"';
        }
        gatestr += '>' + gate.id + ":" + gate.name + '</option>';
      }
      $("#select_gates").html(gatestr);
      $("#select_gates").change(function () {
        window.location.href = 'http://fayaa.com/youxi/hrd/' + this.value + "/";
      });
    });
  }
});

