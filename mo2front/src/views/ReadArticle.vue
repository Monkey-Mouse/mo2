<template>
  <div>
    <aside style="position: fixed; width: 20%" class="mt-16 ml-4">
      <div class="has-icons-right" v-show="!$vuetify.breakpoint.mobile">
        <!-- <a @click="scrollToTop" style="text-decoration: none">{{ title }}</a
        ><v-icon style="float: right">mdi-format-list-bulleted</v-icon
        ><v-divider /> -->
        <div id="toc"></div>
      </div>
    </aside>
    <v-container>
      <v-row id="mo2blog" justify="center">
        <v-col cols="12" lg="7" class="mo2editor">
          <v-skeleton-loader
            v-if="loading"
            v-bind="attrs"
            type="heading, list-item-avatar, paragraph@9"
          ></v-skeleton-loader>
          <div v-if="!loading">
            <div
              id="titleContainer"
              class="mo2title text-break has-icons-right"
            >
              <h1 id="title">{{ title }}</h1>
            </div>
            <v-container v-if="authorLoad">
              <v-row>
                <avatar :size="40" :user="author"></avatar>
                <a
                  @click="$router.push('/account/' + author.id)"
                  class="text--lighten-2 ml-2 mt-2"
                  >{{ author.name }}</a
                >
                <span class="ml-2 grey--text mt-2">{{
                  blog.entityInfo.createTime.substr(0, 10)
                }}</span>
                <v-spacer />
                <v-rating
                  v-if="!draft && blog.score_num"
                  :length="5"
                  :value="blog.score_sum / blog.score_num"
                  readonly
                  half-increments
                ></v-rating>
                <v-tooltip v-if="draft" bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      plain
                      small
                      v-bind="attrs"
                      v-on="on"
                      @click="publish"
                      :style="
                        $vuetify.breakpoint.mobile
                          ? 'padding: 0px;min-width:0px'
                          : ''
                      "
                    >
                      <v-icon>mdi-publish</v-icon>
                    </v-btn>
                  </template>
                  <span>Publish</span>
                </v-tooltip>
                <v-tooltip v-if="user.id === blog.authorId" bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      plain
                      small
                      :style="
                        $vuetify.breakpoint.mobile
                          ? 'padding: 0px;min-width:0px'
                          : ''
                      "
                      v-bind="attrs"
                      v-on="on"
                      @click="edit"
                    >
                      <v-icon>mdi-file-document-edit</v-icon>
                    </v-btn>
                  </template>
                  <span>Edit</span>
                </v-tooltip>
                <v-tooltip v-if="user.id === blog.authorId" bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      @click="deleteArticle"
                      plain
                      small
                      v-bind="attrs"
                      v-on="on"
                      :style="
                        $vuetify.breakpoint.mobile
                          ? 'padding: 0px;min-width:0px'
                          : ''
                      "
                    >
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                  </template>
                  <span>Delete</span>
                </v-tooltip>
                <v-tooltip
                  v-if="blog.entityInfo.is_deleted && user.id === blog.authorId"
                  bottom
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      :style="
                        $vuetify.breakpoint.mobile
                          ? 'padding: 0px;min-width:0px'
                          : ''
                      "
                      @click="restoreArticle"
                      plain
                      small
                      v-bind="attrs"
                      v-on="on"
                    >
                      <v-icon>mdi-delete-restore</v-icon>
                    </v-btn>
                  </template>
                  <span>Restore</span>
                </v-tooltip>
              </v-row>
            </v-container>
            <v-row v-else class="mb-6">
              <v-skeleton-loader
                class="col"
                type="list-item-avatar"
              ></v-skeleton-loader>
            </v-row>
            <!-- <img
          class="ma-5"
          src="https://th.bing.com/th/id/OIP.dnWfZl6P-0Pl47j7PhZodQHaHJ?w=187&h=180&c=7&o=5&dpr=2&pid=1.7"
        />
        <v-row justify="center" class="mb-5">• • •</v-row> -->
            <div
              id="contents"
              v-html="$sanitize(html)"
              class="mo2content mt-10"
              spellcheck="false"
            ></div>
            <delete-confirm
              :title="'确认删除'"
              :content="deleteContent"
              :show.sync="showDelete"
              @confirm="confirmDelete"
            />
            <div style="padding-bottom: 1rem"></div>
            <v-row v-if="!draft" class="pt-6">
              <v-rating
                v-if="!draft"
                :length="5"
                :value="blog.score_sum / blog.score_num"
                @input="rateChange"
                half-increments
              ></v-rating>
              <v-spacer />
              <v-col cols="auto"
                ><v-icon @click="toggleLike">{{
                  liked ? "mdi-thumb-up" : "mdi-thumb-up-outline"
                }}</v-icon
                >{{ praiseNum }}</v-col
              >
              <v-col cols="auto"
                ><v-icon @click="share">mdi-share</v-icon></v-col
              >
              <v-col cols="auto"
                ><v-icon @click="loadComment">mdi-message-reply-outline</v-icon
                >{{ commentNum }}</v-col
              >
            </v-row>

            <!-- 评论 -->
            <div style="padding-bottom: 5rem"></div>
            <v-navigation-drawer
              v-model="comment"
              width="30%"
              height="100%"
              style="max-height: 100%"
              bottom
              fixed
              temporary
            >
              <template v-slot:prepend>
                <v-list-item two-line class="ml-16">
                  <v-list-item-content>
                    <v-icon x-large>mdi-message-reply-outline</v-icon>
                  </v-list-item-content>

                  <v-list-item-content>
                    <v-list-item-title>Comments</v-list-item-title>
                    <!-- <v-list-item-subtitle>Logged In</v-list-item-subtitle> -->
                  </v-list-item-content>
                  <v-list-item-content>
                    <v-icon
                      v-if="$vuetify.breakpoint.mobile"
                      @click="comment = false"
                      x-large
                      >mdi-chevron-triple-down</v-icon
                    >
                    <!-- <v-list-item-subtitle>Logged In</v-list-item-subtitle> -->
                  </v-list-item-content>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item v-if="isUser" class="ma-4"
                  ><v-textarea
                    :loading="commentPosting"
                    auto-grow
                    placeholder="Write what you think about"
                    flat
                    reverse
                    rows="1"
                    v-model="commentmsg"
                    @click="writeCommentShow = true"
                  >
                  </v-textarea>
                  <v-expand-transition>
                    <div v-if="writeCommentShow">
                      <v-icon @click="postComment">mdi-send</v-icon>
                    </div>
                  </v-expand-transition>
                </v-list-item>
              </template>
              <v-skeleton-loader v-if="commentLoading" type="card@3" />
              <v-list v-else v-for="(c, i) in cs" :key="i" nav dense>
                <div>
                  <v-list-item two-line>
                    <v-list-item-avatar class="clickable">
                      <avatar :size="30" :user="c.authorProfile"></avatar>
                    </v-list-item-avatar>

                    <v-list-item-content>
                      <v-list-item-title>{{
                        c.authorProfile.name
                      }}</v-list-item-title>
                      <time-ago
                        :refresh="60"
                        :datetime="c.entity_info.updateTime"
                        tooltip
                        long
                      ></time-ago>
                    </v-list-item-content>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-content>{{ c.content }} </v-list-item-content>
                  </v-list-item>
                  <v-list-item>
                    <v-spacer />
                    <v-icon @click="loadSub(c)"
                      >mdi-message-reply-outline</v-icon
                    >{{ c.subs.length }}
                    <v-list-item-action
                      ><v-icon @click="c.edit = !c.edit"
                        >mdi-reply</v-icon
                      ></v-list-item-action
                    >
                  </v-list-item>
                  <v-expand-transition>
                    <v-list-item v-if="c.edit" class="ma-4"
                      ><v-textarea
                        :loading="commentPosting"
                        auto-grow
                        placeholder="Write what you think about"
                        flat
                        reverse
                        rows="1"
                        v-model="c.tempC"
                      >
                      </v-textarea>
                      <div>
                        <v-icon @click="postSubComment(c)">mdi-send</v-icon>
                      </div>
                    </v-list-item>
                  </v-expand-transition>
                  <v-divider />
                  <div v-if="c.showSub">
                    <v-list
                      class="ml-16"
                      v-for="(s, i) in c.subs"
                      :key="i"
                      nav
                      dense
                    >
                      <div>
                        <v-list-item two-line>
                          <v-list-item-avatar class="clickable">
                            <avatar :size="30" :user="s.authorProfile"></avatar>
                          </v-list-item-avatar>

                          <v-list-item-content>
                            <v-list-item-title>{{
                              s.authorProfile.name
                            }}</v-list-item-title>
                            <time-ago
                              :refresh="60"
                              :datetime="s.entity_info.updateTime"
                              tooltip
                              long
                            ></time-ago>
                          </v-list-item-content>
                        </v-list-item>
                        <v-list-item>
                          <v-list-item-content
                            >{{ s.content }}
                          </v-list-item-content>
                        </v-list-item>
                      </div>
                      <v-divider />
                    </v-list>
                  </div>
                </div>
              </v-list>
              <v-skeleton-loader v-if="commentLoadingMore" type="card@3" />
              <v-list v-if="!nomore">
                <v-row justify="center" class="text-center">
                  <v-btn
                    @click="loadMoreComments"
                    class="ma-5"
                    fab
                    dark
                    color="primary"
                  >
                    <v-icon dark> mdi-plus </v-icon>
                  </v-btn></v-row
                ></v-list
              >
              <v-list v-if="nomore && commentNum === 0">
                <v-list-item>
                  <h1 class="ml-7">暂时没有评论</h1>
                </v-list-item>
              </v-list>
            </v-navigation-drawer>
            <MO2Dialog
              :show.sync="showPublish"
              confirmText="发布"
              title="发布文章"
              :inputProps="inputProps"
              :validator="validator"
              ref="dialog"
              :confirm="confirm"
              :uploadImgs="uploadImgs"
            />
          </div>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Editor from "../components/Editor/MO2Editor.vue";
import {
  DeleteArticle,
  GenerateTOC,
  GetArticle,
  GetBlogLikeNum,
  GetCommentNum,
  GetComments,
  GetErrorMsg,
  GetUserData,
  GetUserDatas,
  globaldic,
  InitLoader,
  IsBlogLiked,
  IsScoredBlog,
  Prompt,
  RecycleBlog,
  RestoreBlog,
  ScoreBlog,
  ShareToQQ,
  ShowLogin,
  ToggleLikeBlog,
  UploadImgToQiniu,
  UpsertBlog,
  UpsertComment,
  UpsertSubComment,
  UserRole,
} from "@/utils";
import hljs from "highlight.js";
import {
  Blog,
  User,
  Comment,
  UserListData,
  BlankBlog,
  InputProp,
  BlogUpsert,
} from "@/models";
import Avatar from "@/components/UserAvatar.vue";
import { Prop, Watch } from "vue-property-decorator";
import { TimeAgo } from "vue2-timeago";
import DeleteConfirm from "@/components/DeleteConfirm.vue";
import { required, minLength } from "vuelidate/lib/validators";
import MO2Dialog from "../components/MO2Dialog.vue";
@Component({
  components: {
    Editor,
    Avatar,
    DeleteConfirm,
    TimeAgo,
    MO2Dialog,
  },
})
export default class ReadArticle extends Vue {
  @Prop()
  user: User;
  uploadImgs = UploadImgToQiniu;
  showPublish = false;
  validator = {
    description: {
      required: required,
      min: minLength(8),
    },
    title: {
      required: required,
    },
  };
  inputProps: { [name: string]: InputProp } = {
    title: {
      errorMsg: {
        required: "标题不可为空",
      },
      label: "Title",
      default: "",
      icon: "mdi-format-title",
      col: 12,
      type: "text",
    },
    description: {
      errorMsg: {
        required: "描述不可为空",
        min: "描述长度不小于8",
      },
      label: "Description",
      default: "",
      icon: "mdi-text",
      col: 12,
      type: "textarea",
    },
    cover: {
      errorMsg: {},
      label: "Cover",
      default: {},
      icon: "mdi-image",
      col: 12,
      type: "imgselector",
    },
    categories: {
      errorMsg: {},
      label: "Categories",
      default: "",
      col: 12,
      options: [],
      type: "select",
    },
  };
  get dialog() {
    return this.$refs["dialog"] as MO2Dialog;
  }
  title = "";
  html = "";
  writeCommentShow = false;
  attrs = {
    class: "mb-6 mt-6",
    boilerplate: false,
    elevation: 0,
  };
  loading = true;
  blog: Blog = BlankBlog;
  author: User;
  authorLoad = false;
  showDelete = false;
  draft = false;
  comment = false;
  commentmsg = "";
  p = 0;
  ps = 5;
  cs: Comment[] = [];
  commentLoading = true;
  commentPosting = false;
  commentLoadingMore = true;
  nomore = false;
  commentNum = 0;
  praiseNum = 0;
  liked = false;

  get isUser() {
    return this.user.roles && this.user.roles.indexOf(UserRole) >= 0;
  }
  async rateChange(rate: number) {
    const p = IsScoredBlog({ score: rate, target: this.blog.id });
    const re = await ScoreBlog({ score: rate, target: this.blog.id });
    this.blog.score_sum = re.sum;
    this.blog.score_num = re.num;
    const sed = await p;
    Prompt((sed ? "重新" : "") + "打分成功！", 3000);
  }
  get deleteContent() {
    return (
      "你确定要删除" +
      (this.draft ? "草稿" : "") +
      '"' +
      this.title +
      '"吗？' +
      (this.blog.entityInfo.is_deleted || this.draft
        ? "此文章将被彻底删除，这个操作不能被撤销"
        : "删除后此文章将保存在垃圾箱里24小时")
    );
  }
  async confirm(model: BlogUpsert, draft = false) {
    try {
      for (const key in model) {
        if (Object.prototype.hasOwnProperty.call(model, key)) {
          const element = model[key];
          this.blog[key] = element;
        }
      }
      await UpsertBlog({ draft: false }, this.blog);
      this.$router.push("/article/" + this.blog.id);
      return { err: "", pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  restoreArticle() {
    RestoreBlog(this.blog.id, { draft: this.draft })
      .then(() => {
        this.blog.entityInfo.is_deleted = false;
        Prompt("成功恢复文章！", 10000);
      })
      .catch((err) => {
        Prompt(GetErrorMsg(err), 10000);
      });
  }
  publish() {
    this.showPublish = true;
    console.log(this.blog);
    this.dialog.setModel(this.blog);
    let imgEs = document.querySelectorAll(".mo2content img");
    const array = [...imgEs];
    const list = [];
    if (this.blog.cover && this.blog.cover !== "") {
      list.push({
        src: this.blog.cover,
        active: false,
      });
    }
    array.map((e) => {
      const i = e as HTMLImageElement;
      list.push({ src: i.src, active: false });
    });
    list.push({
      src: "//cdn.mo2.leezeeyee.com/60365aae06fd3124561400c3/1614260703850314778image.png",
      active: false,
    });
    list[0].active = true;
    this.dialog.setModel({ cover: list });
  }
  scrollToTop() {
    window.scrollTo(0, 0);
  }
  toggleLike() {
    if (this.isUser) {
      ToggleLikeBlog(this.blog.id).then(() => {
        if (this.liked) {
          this.praiseNum--;
        } else this.praiseNum++;
        this.liked = !this.liked;
      });
    } else {
      Prompt("请先登录！", 5000);
      ShowLogin();
    }
  }
  share() {
    ShareToQQ({
      title: this.blog.title,
      pic: this.blog.cover,
      summary: this.blog.description,
      desc: this.blog.description,
    });
  }
  @Watch("$route")
  articleChange() {
    if (this.$route.params["id"] !== this.blog.id) {
      this.init();
    }
  }
  created() {
    this.init();
  }
  init() {
    const toc = document.getElementById("toc");
    if (toc) {
      toc.innerHTML = "";
    }

    this.loading = true;
    if (this.$route.query["draft"]) {
      this.draft = (this.$route.query["draft"] as string) === "true";
    }
    GetArticle({ id: this.$route.params["id"], draft: this.draft })
      .then((val) => {
        this.loading = false;
        this.blog = val;
        this.title = val.title;
        this.html = val.content;
        GetUserData(this.blog.authorId).then((u) => {
          this.author = u;
          this.authorLoad = true;
        });
        GetCommentNum(this.blog.id).then((c) => {
          this.commentNum = c.count;
        });
        GetBlogLikeNum(this.blog.id).then((num) => {
          this.praiseNum = num.num;
        });
        if (this.isUser) {
          IsBlogLiked(this.blog.id).then((l) => {
            this.liked = l.liked;
          });
        }

        setTimeout(() => {
          GenerateTOC();
          // first, find all the code blocks
          document.querySelectorAll("code").forEach((block) => {
            // then highlight each
            hljs.highlightBlock(block);
          });
        }, 50);
      })
      .catch((err) => GetErrorMsg(err));
  }
  async loadSub(c: Comment) {
    if (c.showSub === true) {
      c.showSub = false;
      return;
    }
    if (c.subs.length > 0 && c.subs[0].authorProfile) {
      c.showSub = true;
      return;
    }
    const map: { [key: string]: UserListData } = {};
    (await GetUserDatas(c.subs.map((v) => v.aurhor))).forEach(
      (v) => (map[v.id] = v)
    );
    c.subs.forEach((v) => {
      v.authorProfile = map[v.aurhor];
      v.edit = false;
    });
    c.showSub = true;
  }
  async postSubComment(c: Comment) {
    this.commentPosting = true;
    const sub = await UpsertSubComment(c.id, { content: c.tempC });
    sub.authorProfile = this.user;
    c.showSub = true;
    c.subs.unshift(sub);
    c.tempC = "";
    this.commentPosting = false;
  }
  async postComment() {
    this.commentPosting = true;
    const c = await UpsertComment({
      article: this.blog.id,
      content: this.commentmsg,
    });
    c.authorProfile = this.user;
    c.edit = false;
    c.tempC = "";
    c.showSub = false;
    this.cs.unshift(c);
    this.commentmsg = "";
    this.commentNum++;
    this.commentPosting = false;
  }
  edit() {
    globaldic.article = `<h1>${this.title}</h1>${this.html}`;
    // UpsertBlog({ draft: true }, this.blog);
    this.$router.push("/edit/" + this.blog.id);
  }
  async loadComment() {
    this.comment = true;
    await this.loadMoreComments();
    this.commentLoading = false;
  }
  async loadMoreComments() {
    this.commentLoadingMore = true;
    const newCs = await GetComments(this.blog.id, {
      page: this.p++,
      pagesize: this.ps,
    });
    if (newCs.length < this.ps) {
      this.nomore = true;
    }
    const map: { [key: string]: UserListData } = {};
    (await GetUserDatas(newCs.map((v) => v.aurhor))).forEach(
      (v) => (map[v.id] = v)
    );
    newCs.forEach((v) => {
      v.authorProfile = map[v.aurhor];
      v.edit = false;
      v.tempC = "";
      v.showSub = false;
    });
    this.cs = this.cs.concat(newCs);
    this.commentLoadingMore = false;
  }
  deleteArticle() {
    this.showDelete = true;
  }
  confirmDelete() {
    let delFunc = RecycleBlog;
    if (this.blog.entityInfo.is_deleted || this.draft) {
      delFunc = DeleteArticle;
    }
    delFunc(this.blog.id, { draft: this.draft })
      .then(() => {
        Prompt("成功删除文章！", 10000);
      })
      .catch((err) => {
        Prompt(GetErrorMsg(err), 10000);
      });
  }
}
</script>
<style>
.v-skeleton-loader__text {
  border-radius: 6px;
  flex: 1 0 auto;
  height: 12px;
  margin-bottom: 6px !important;
}
</style>
<style scoped>
@import url("https://cdn.jsdelivr.net/npm/katex@0.13.5/dist/katex.min.css");
</style>
<style lang="scss">
.v-application {
  span.katex {
    .accent {
      background-color: unset !important;
      border-color: unset !important;
    }
  }
}
</style>