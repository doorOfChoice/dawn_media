function replyOnclick(e) {
    var box = $("#comment-box");
    var now = $(e);
    if(now.attr("reply") !== "true") {
        box.find('#parent_id').val(now.attr("data-id"));
        now.before(box);
        now.attr("reply", "true")
        now.text("收起");
    }else {
        $(".comment").append(box);
        now.attr("reply", "");
        now.text("回复");
    }
}
$(function () {
    var comments = $("#comments");
    var media_id = $("#media_id").val();
    var read_more = $("#read-more-comments");
    var pageTpl = "<p>总评:<%:=page.Count%></p>";
    var tpl = `
        <%for(var i = 0; i < comments.length; i++){%>
        <div class="row">
        <div class="col-xs-12 col-md-2">
            <img class="user-head-medium" src="<%:=avatarMap%><%:=comments[i].User.Avatar%>" alt="">
            <center><span><%:=comments[i].User.Nickname%></span></center>
        </div>
        <div class="col-xs-12 col-md-10">
            <p><%:=comments[i].CreatedAt%></p>
            <%if(comments[i].ParentComment != null){%>
                <div class="alert alert-danger" role="alert">
                   <p>
                   <b>回复@<%:=comments[i].ParentComment.User.Nickname%></b>
                   :
                   <%:=comments[i].ParentComment.Content%>
                   </p>
                </div>
            <%}%>
            </p>
            <%:=comments[i].Content%>
            </p>
            <p><b><span data-id="<%:=comments[i].ID%>" onclick="replyOnclick(this)">回复</span></b></p>
        </div>
        </div>
        <div class="line"></div>
        <%}%>
    `;
    template.registerFunction("getParentUser", function (parentComment) {
        if(parentComment == null)
            return "";

        return "@" + parentComment.User.Nickname + ": " + parentComment.Content;
    });



    function load(uri) {
        axios.get(uri)
            .then(function (rep) {
                var cs = rep.data.data.comments;
                var page = rep.data.data.page;
                for(var i = 0; i < cs.length; i++) {
                    cs[i].CreatedAt = moment(cs[i].CreatedAt).format("YYYY年M月D  H:mm:ss");
                }
                var result = template(pageTpl + tpl, {
                    comments :cs,
                    avatarMap : rep.data.data.avatarMap,
                    page:page,
                });
                comments.append($(result));
                if(page.CurPage < page.MaxPage) {
                    read_more.attr("next-link", page.NextLink);
                }else {
                    read_more.attr("next-link", "#");
                    read_more.text("没有更多了...");
                }
            })
    }
    read_more.click(function (e) {
        var uri = $(this).attr("next-link");
        if(uri !== "#") {
            load(uri);
        }
    });
    $("#comment-submit").click(function (e) {
        var box = $("#comment-box");
        axios.post("/comments", {
            "Parent_ID" : $("#parent_id").val(),
            "Media_ID": $("#media_id").val(),
            "Content":$("#comment-content").val()
        }).then(function (rep) {
            var comment = rep.data.data.comment;
            comment.CreatedAt = moment(comment.CreatedAt).format("YYYY年M月D  H:mm:ss");
            var result = template(tpl, {
                comments :[comment],
                avatarMap : rep.data.data.avatarMap
            });
            comments.prepend($(result));
        }).catch(function (err) {
            alert(err.response.data.error);
        });
    });

    $("#star-btn").click(function (e) {
        var that = this;
       axios.post("/media/star", {
           "media_id" : media_id
       }).then(function (rep) {
            var data = rep.data.data;
            $(that).find(".star-count").text(data.count);
            var heart = $(that).find("i");
            if(data.create) {
                heart.addClass("red");
            }else {
                heart.removeClass("red");
            }
       }).catch(function (err) {
           alert(err.response.data.error);
       });
    });

    load( "/comments?id=" + media_id + "&" + "curTime=" + moment().unix());
    var player = videojs('video');
    player.videoJsResolutionSwitcher();
    player.on('error', function (e) {
       console.log(e);
       player.error(null);
       player.play();
    });

});