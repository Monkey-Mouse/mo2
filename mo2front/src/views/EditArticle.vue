<template>
  <div>
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <v-skeleton-loader
          class="mb-10 mt-10"
          v-if="loading"
          type="heading"
        ></v-skeleton-loader>
        <v-skeleton-loader
          v-if="loading"
          type="paragraph@9"
        ></v-skeleton-loader>
      </v-col>
    </v-row>
    <editor
      v-show="!loading"
      :uploadImgs="uploadImgs"
      :content="content"
      @loaded="editorLoad"
    />
    <MO2Dialog
      :show.sync="showPublish"
      confirmText="发布"
      title="发布文章"
      :inputProps="inputProps"
      :validator="validator"
      ref="dialog"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Editor from "../components/MO2Editor.vue";
import MO2Dialog from "../components/MO2Dialog.vue";
import { globaldic, UploadImgToQiniu } from "@/utils";
import { BlogUpsert, InputProp } from "@/models";
import { required, minLength, email } from "vuelidate/lib/validators";
@Component({
  components: {
    Editor,
    MO2Dialog,
  },
})
export default class EditArticle extends Vue {
  showPublish = false;
  uploadImgs = UploadImgToQiniu;
  loading = true;
  editorLoaded = false;
  content = "";
  editor: Editor;
  blog: BlogUpsert = {};
  validator = {
    description: {
      required: required,
      min: minLength(8),
    },
    title: {
      required: required,
    },
    name: {
      required: required,
    },
  };
  get inputProps(): { [name: string]: InputProp } {
    return {
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
        },
        label: "Description",
        default: "",
        icon: "mdi-text",
        col: 12,
        type: "textarea",
      },
    };
  }
  editorLoad(editor: Editor) {
    this.editor = editor;
    this.editorLoaded = true;
    if (!this.$route.params["id"] || this.content !== "") {
      this.loading = false;
    }
  }
  publish() {
    const titleElm = document.querySelector("h1");
    this.blog.title = titleElm.innerText.trim();
    this.blog.content = this.editor
      .GetHTML()
      .substring(9 + titleElm.innerText.length);
    let elm = document.querySelector(".mo2content p") as any;
    this.blog.description = "";
    const descriptions = [];
    while (this.blog.description.length < 50 && elm) {
      descriptions.push(elm.innerText);
      while (elm.nextElementSibling && elm.nextElementSibling.tagName !== "P") {
        elm = elm.nextElementSibling;
      }
      elm = elm.nextElementSibling;
    }

    this.blog.description = descriptions.join("").trim();
    this.showPublish = true;
    (this.$refs["dialog"] as MO2Dialog).setModel(this.blog);
  }
}
</script>
