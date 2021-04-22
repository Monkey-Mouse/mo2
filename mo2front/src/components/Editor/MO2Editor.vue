<template>
  <v-container>
    <input
      type="file"
      ref="f"
      accept="image/*"
      style="display: none"
      @change="fileSelected"
    />
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <!-- <editor-menu-bar :editor="editor" v-slot="{ commands, isActive }">
          <div class="menubar" style="z-index: 9999">
            <div class="toolbar row">
              <span v-if="isActive.table()">
                <button class="menubar__button" @click="commands.deleteTable">
                  <v-icon>mdi-table-large-remove</v-icon>
                </button>
                <button
                  class="menubar__button"
                  @click="commands.addColumnBefore"
                >
                  <v-icon>mdi-table-column-plus-before</v-icon>
                </button>
                <button
                  class="menubar__button"
                  @click="commands.addColumnAfter"
                >
                  <v-icon>mdi-table-column-plus-after</v-icon>
                </button>
                <button class="menubar__button" @click="commands.deleteColumn">
                  <v-icon>mdi-table-column-remove</v-icon>
                </button>
                <button class="menubar__button" @click="commands.addRowBefore">
                  <v-icon>mdi-table-row-plus-before</v-icon>
                </button>
                <button class="menubar__button" @click="commands.addRowAfter">
                  <v-icon>mdi-table-row-plus-after</v-icon>
                </button>
                <button class="menubar__button" @click="commands.deleteRow">
                  <v-icon>mdi-table-row-remove</v-icon>
                </button>
              </span>
            </div>
          </div>
        </editor-menu-bar> -->
        <v-progress-circular
          v-if="isuploading"
          class="offset-10 offset-lg-7"
          style="position: fixed; z-index: 999"
          indeterminate
          color="amber"
        ></v-progress-circular>
        <bubble-menu v-if="editor" :editor="editor">
          <v-btn-toggle dense :value="[]">
            <v-btn
              @click="editor.chain().focus().toggleBold().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('bold'),
              }"
            >
              <v-icon>mdi-format-bold</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().toggleItalic().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('italic'),
              }"
            >
              <v-icon>mdi-format-italic</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().toggleStrike().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('strike'),
              }"
            >
              <v-icon>mdi-format-strikethrough</v-icon>
            </v-btn>
          </v-btn-toggle>
        </bubble-menu>
        <floating-menu v-if="editor" :editor="editor">
          <v-btn-toggle dense :value="[]">
            <v-btn
              small
              @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('heading', {
                  level: 1,
                }),
              }"
            >
              <v-icon>mdi-format-header-1</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('heading', {
                  level: 2,
                }),
              }"
            >
              <v-icon>mdi-format-header-2</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleBulletList().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('bulletList'),
              }"
            >
              <v-icon>mdi-format-list-bulleted</v-icon>
            </v-btn>
          </v-btn-toggle>
        </floating-menu>
        <div class="mo2content" spellcheck="false">
          <editor-content v-if="editor" :editor="editor" />
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {
  Editor,
  EditorContent,
  FloatingMenu,
  BubbleMenu,
  VueNodeViewRenderer,
} from "@tiptap/vue-2";
import Placeholder from "@tiptap/extension-placeholder";
import Paragraph from "@tiptap/extension-paragraph";
import Text from "@tiptap/extension-text";
import BlockQuote from "@tiptap/extension-blockquote";
import Italic from "@tiptap/extension-italic";
import Bold from "@tiptap/extension-bold";
import Strike from "@tiptap/extension-strike";
import Heading from "@tiptap/extension-heading";
import BulletList from "@tiptap/extension-bullet-list";
import OrderedList from "@tiptap/extension-ordered-list";
import ListItem from "@tiptap/extension-list-item";
import Doc from "@tiptap/extension-document";
import { Prop, Watch } from "vue-property-decorator";
import { timeout } from "@/utils";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import CodeBlockComponent from "./Lowlight/CodeBlockComponent.vue";
// load all highlight.js languages
import lowlight from "lowlight";

let that: MO2Editor | any = {};
@Component({
  components: {
    EditorContent,
    FloatingMenu,
    BubbleMenu,
    // EditorMenuBar,
  },
})
export default class MO2Editor extends Vue {
  @Prop()
  content?: string;
  @Prop()
  uploadImgs: (
    blobs: File[],
    callback: (imgprop: { src: string }) => void
  ) => Promise<void>;
  isuploading = false;
  linkMenuIsActive = false;
  update = false;
  editable = true;
  load = false;
  editor: Editor = null;
  initEditor(content: string) {
    console.log("init");
    if (this.editor) {
      this.editor.commands.setContent(content ?? "<h1></h1>");
      return;
    }
    if (!content || content.length === 0) {
      content = "<h1></h1><p></p>";
    }
    this.editor = new Editor({
      extensions: [
        Paragraph,
        Text,
        Doc,
        Italic,
        Bold,
        BulletList,
        OrderedList,
        Strike,
        BlockQuote,
        ListItem,
        Placeholder.configure({ showOnlyCurrent: false }),
        Heading.configure({ levels: [1, 2, 3, 4] }),
        CodeBlockLowlight.extend({
          addNodeView() {
            return VueNodeViewRenderer(CodeBlockComponent as any);
          },
        }).configure({ lowlight }),
      ],
      content: content ?? "<h1></h1>",
      onUpdate() {
        (that as MO2Editor).update = true;
      },
    });
  }
  fileSelected() {
    this.isuploading = true;
    // this.uploadImgs(
    //   [...(this.$refs.f as HTMLInputElement).files],
    //   this.editor.commands.image
    // )
    //   .then(() => (that.isuploading = false))
    //   .catch(() => {
    //     that.isuploading = false;
    //   });
    (this.$refs.f as HTMLInputElement).value = "";
  }
  mounted() {
    that = this;
    this.initEditor(this.content);
    this.$emit("loaded", this);
    this.startAutoSave();
  }
  async startAutoSave() {
    while (true) {
      if (this.update) {
        (that as MO2Editor).$emit("autosave");
        this.update = false;
      }
      await timeout(5000);
    }
  }

  @Watch("content")
  contentSet() {
    if (this.content !== null && this.content !== undefined) {
      this.initEditor(this.content);
    }
  }

  beforeDestroy() {
    this.editor.destroy();
  }
  linkUrl = null;
  showLinkMenu(attrs) {
    this.linkUrl = attrs.href;
    this.linkMenuIsActive = true;
    this.$nextTick(() => {
      (this.$refs.linkInput as any).focus();
    });
  }
  hideLinkMenu() {
    this.linkUrl = null;
    this.linkMenuIsActive = false;
  }
  setLinkUrl(command, url) {
    command({ href: url });
    this.hideLinkMenu();
  }
  @Watch("editable")
  changeEditable() {
    this.editor.setOptions({
      editable: this.editable,
    });
  }
  public GetHTML() {
    return this.editor.getHTML() as string;
  }
}
</script>