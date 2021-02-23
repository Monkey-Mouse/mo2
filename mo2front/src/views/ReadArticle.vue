<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <v-skeleton-loader
          v-if="loading"
          v-bind="attrs"
          type="heading, list-item-avatar, paragraph@9"
        ></v-skeleton-loader>
        <div v-show="!loading">
          <div class="mo2title">
            <h1>{{ title }}</h1>
          </div>
          <v-row class="mb-6">
            <v-col>
              <v-avatar color="primary" size="40">me</v-avatar>
              <span class="text--lighten-2 ml-2">李子怡</span>
              <span class="ml-2 grey--text">2020-1-1</span>
            </v-col>
            <v-spacer />
            <v-btn plain small>
              <v-icon @click="edit">mdi-file-document-edit</v-icon>
            </v-btn>
            <v-btn plain small>
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </v-row>
          <!-- <img
          class="ma-5"
          src="https://th.bing.com/th/id/OIP.dnWfZl6P-0Pl47j7PhZodQHaHJ?w=187&h=180&c=7&o=5&dpr=2&pid=1.7"
        />
        <v-row justify="center" class="mb-5">• • •</v-row> -->
          <div v-html="html" class="mo2content mt-10" spellcheck="false"></div>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Editor from "../components/MO2Editor.vue";
import { GetArticle, globaldic, UploadImgToQiniu } from "@/utils";
import hljs from "highlight.js";
import { Blog } from "@/models";
@Component({
  components: {
    Editor,
  },
})
export default class ReadArticle extends Vue {
  title = "";
  html = "";
  attrs = {
    class: "mb-6 mt-6",
    boilerplate: false,
    elevation: 0,
  };
  loading = true;
  blog: Blog;
  created() {
    var draft = false;
    if (this.$route.query["draft"]) {
      draft = (this.$route.query["draft"] as string) === "true";
    }
    GetArticle({ id: this.$route.params["id"], draft: draft }).then((val) => {
      this.loading = false;
      this.blog = val;
      this.title = val.title;
      this.html = val.content;
      setTimeout(() => {
        // first, find all the code blocks
        document.querySelectorAll("code").forEach((block) => {
          // then highlight each
          hljs.highlightBlock(block);
        });
      }, 50);
    });
  }
  edit() {
    globaldic.article = `<h1>${this.title}</h1>${this.html}`;
    this.$router.push("/edit/" + this.blog.id);
  }
}
</script>
