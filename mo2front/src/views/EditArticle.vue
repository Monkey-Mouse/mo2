<template>
  <div>
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <v-skeleton-loader
          class="mb-10 mt-10"
          v-if="loading || !editorLoaded"
          type="heading"
        ></v-skeleton-loader>
        <v-skeleton-loader
          v-if="loading || !editorLoaded"
          type="paragraph@9"
        ></v-skeleton-loader>
      </v-col>
    </v-row>
    <editor
      v-show="!(loading || !editorLoaded)"
      :uploadImgs="uploadImgs"
      :content="content"
      @loaded="editorLoad"
      @autosave="autoSave"
    />
    <MO2Dialog
      :show.sync="showPublish"
      confirmText="发布"
      title="发布文章"
      :inputProps="inputProps"
      :validator="validator"
      ref="dialog"
      @confirm="confirm"
      :uploadImgs="uploadImgs"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Editor from "../components/MO2Editor.vue";
import MO2Dialog from "../components/MO2Dialog.vue";
import {
  GetArticle,
  GetErrorMsg,
  globaldic,
  UploadImgToQiniu,
  UpsertBlog,
} from "@/utils";
import { BlogUpsert, InputProp } from "@/models";
import { required, minLength, email } from "vuelidate/lib/validators";
import { AxiosError } from "axios";
import { Prop } from "vue-property-decorator";
@Component({
  components: {
    Editor,
    MO2Dialog,
  },
})
export default class EditArticle extends Vue {
  @Prop()
  autoSaving: boolean;
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
      cover: {
        errorMsg: {},
        label: "Cover",
        default: {},
        icon: "mdi-image",
        col: 12,
        type: "imgselector",
      },
    };
  }
  created() {
    // this.content = globaldic.article ?? "";
    // globaldic.article = "";
    // console.log(this.content, this.$route.params["id"]);
    if (this.content !== "") {
      this.loading = false;
      this.blog.id = this.$route.params["id"];
    } else if (this.$route.params["id"]) {
      this.blog.id = this.$route.params["id"];
      GetArticle({ id: this.blog.id, draft: true })
        .then((val) => {
          this.blog = val;
          this.content = val.content;
          this.loading = false;
        })
        .catch((reason: AxiosError) => {
          if (reason.response.status === 404) {
            GetArticle({ id: this.blog.id, draft: false })
              .then((val) => {
                this.blog = val;
                this.content = val.content;
                this.loading = false;
              })
              .catch((err) => {});
          }
        });
    }
    window.addEventListener("beforeunload", () => {
      this.getTitleAndContent();
      this.postBlog({}, true);
    });
  }
  editorLoad(editor: Editor) {
    this.editor = editor;
    this.editorLoaded = true;
  }
  getTitleAndContent() {
    const titleElm = document.querySelector("h1");
    this.blog.title = titleElm.innerText.trim();
    this.blog.content = this.editor
      .GetHTML()
      .substring(9 + titleElm.innerText.length);
  }
  publish() {
    this.getTitleAndContent();
    let elm = document.querySelector(".mo2content p") as any;
    if (this.blog.description === "" || this.blog.description === undefined) {
      this.blog.description = "";
      const descriptions = [];
      while (this.blog.description.length < 50 && elm) {
        descriptions.push(elm.innerText);
        while (
          elm.nextElementSibling &&
          elm.nextElementSibling.tagName !== "P"
        ) {
          elm = elm.nextElementSibling;
        }
        elm = elm.nextElementSibling;
      }

      this.blog.description = descriptions.join("").trim();
    }
    this.showPublish = true;
    (this.$refs["dialog"] as MO2Dialog).setModel(this.blog);
    let imgEs = document.querySelectorAll(".mo2content img");
    const array = [...imgEs];
    const list = array.map((e) => {
      const i = e as HTMLImageElement;
      return { src: i.src, active: false };
    });
    list[0].active = true;
    (this.$refs["dialog"] as MO2Dialog).setModel({ cover: list });
  }
  async postBlog(model: BlogUpsert, draft = false) {
    for (const key in model) {
      const element = model[key];
      this.blog[key] = element;
    }
    var data = await UpsertBlog({ draft: draft }, this.blog);
    this.blog.id = data.id;
  }
  async confirm(model: BlogUpsert, draft = false) {
    await this.postBlog(model, draft);
    this.$router.push("/article/" + this.blog.id);
  }
  autoSave() {
    this.$emit("update:autoSaving", true);
    this.getTitleAndContent();
    this.postBlog({}, true)
      .then(() => {
        this.$emit("update:autoSaving", false);
      })
      .catch(() => {
        this.$emit("update:autoSaving", null);
      });
  }
}
</script>
