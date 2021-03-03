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
      v-show="!loading && editorLoaded"
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
  UpSertBlogSync,
} from "@/utils";
import { BlogUpsert, InputProp } from "@/models";
import { required, minLength, email } from "vuelidate/lib/validators";
import { AxiosError } from "axios";
import { Prop, Watch } from "vue-property-decorator";
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
    this.init();
  }
  @Watch("$route", { immediate: true, deep: true })
  pageChange() {
    if (!this.$route.params["id"] || this.$route.params["id"] === "") {
      this.content = "";
      this.blog = {};
      this.init();
    }
  }
  init() {
    if (this.$route.params["id"]) {
      this.blog.id = this.$route.params["id"];
      GetArticle({ id: this.blog.id, draft: true })
        .then((val) => {
          this.blog = val;
          this.content = `<h1>${val.title}</h1>${val.content}`;
          this.loading = false;
        })
        .catch((reason: AxiosError) => {
          if (reason.response.status === 404) {
            GetArticle({ id: this.blog.id, draft: false })
              .then((val) => {
                this.blog = val;
                this.content = `<h1>${val.title}</h1>${val.content}`;
                this.loading = false;
              })
              .catch((err) => {});
          }
        });
    } else this.loading = false;
    window.addEventListener("beforeunload", () => {
      this.getTitleAndContent();
      UpSertBlogSync({ draft: true }, this.blog);
    });
  }
  beforeDestroy() {
    this.getTitleAndContent();
    UpSertBlogSync({ draft: true }, this.blog);
  }
  editorLoad(editor: Editor) {
    this.editor = editor;
    this.editorLoaded = true;
  }
  getTitleAndContent() {
    const raw = this.editor.GetHTML();
    const titlePos = raw.indexOf("</h1>");
    this.blog.title = raw.substring(4, titlePos);
    this.blog.content = raw.substring(titlePos + 5);
  }
  publish() {
    this.getTitleAndContent();
    let elm = document.querySelector(".mo2content p") as any;
    if (this.blog.description === "" || this.blog.description === undefined) {
      this.blog.description = "";
      while (this.blog.description.length < 50 && elm) {
        this.blog.description += (elm.innerText as string).trim();
        while (
          elm.nextElementSibling &&
          elm.nextElementSibling.tagName !== "P"
        ) {
          elm = elm.nextElementSibling;
        }
        elm = elm.nextElementSibling;
      }
    }
    this.showPublish = true;
    (this.$refs["dialog"] as MO2Dialog).setModel(this.blog);
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
      src:
        "//cdn.mo2.leezeeyee.com/60365aae06fd3124561400c3/1614260703850314778image.png",
      active: false,
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
    if (this.blog.id !== data.id) {
      this.$router.replace(`/edit/${data.id}`);
    }
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
