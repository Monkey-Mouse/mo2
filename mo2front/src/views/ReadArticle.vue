<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <v-skeleton-loader
          v-if="loading"
          v-bind="attrs"
          type="heading, list-item-avatar, paragraph@9"
        ></v-skeleton-loader>
        <div v-if="!loading">
          <div class="mo2title text-break">
            <h1>{{ title }}</h1>
          </div>
          <v-row v-if="authorLoad" class="mb-6">
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
            <v-tooltip v-if="draft" bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-btn plain small v-bind="attrs" v-on="on">
                  <v-icon>mdi-eye-check</v-icon>
                </v-btn>
              </template>
              <span>This is a draft</span>
            </v-tooltip>
            <v-btn @click="edit" v-if="user.id === blog.authorId" plain small>
              <v-icon>mdi-file-document-edit</v-icon>
            </v-btn>
            <v-btn
              @click="deleteArticle"
              v-if="user.id === blog.authorId"
              plain
              small
            >
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </v-row>
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
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Editor from "../components/MO2Editor.vue";
import {
  DeleteArticle,
  GetArticle,
  GetErrorMsg,
  GetUserData,
  globaldic,
  UploadImgToQiniu,
  UpsertBlog,
} from "@/utils";
import hljs from "highlight.js";
import { Blog, User } from "@/models";
import Avatar from "@/components/UserAvatar.vue";
import { Prop } from "vue-property-decorator";
import DeleteConfirm from "@/components/DeleteConfirm.vue";
@Component({
  components: {
    Editor,
    Avatar,
    DeleteConfirm,
  },
})
export default class ReadArticle extends Vue {
  @Prop()
  user;
  title = "";
  html = "";
  attrs = {
    class: "mb-6 mt-6",
    boilerplate: false,
    elevation: 0,
  };
  loading = true;
  blog: Blog;
  author: User;
  authorLoad = false;
  showDelete = false;
  draft = false;
  get deleteContent() {
    return '你确定要删除"' + this.title + '"吗？';
  }
  created() {
    if (this.$route.query["draft"]) {
      this.draft = (this.$route.query["draft"] as string) === "true";
    }
    GetArticle({ id: this.$route.params["id"], draft: this.draft }).then(
      (val) => {
        this.loading = false;
        this.blog = val;
        this.title = val.title;
        this.html = val.content;
        GetUserData(this.blog.authorId).then((u) => {
          this.author = u;
          this.authorLoad = true;
        });
        setTimeout(() => {
          // first, find all the code blocks
          document.querySelectorAll("code").forEach((block) => {
            // then highlight each
            hljs.highlightBlock(block);
          });
        }, 50);
      }
    );
  }
  edit() {
    globaldic.article = `<h1>${this.title}</h1>${this.html}`;
    // UpsertBlog({ draft: true }, this.blog);
    this.$router.push("/edit/" + this.blog.id);
  }
  deleteArticle() {
    this.showDelete = true;
  }
  confirmDelete() {
    DeleteArticle(this.blog.id, { draft: this.draft })
      .then(() => {
        this.$router.back();
      })
      .catch((err) => {
        alert(GetErrorMsg(err));
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