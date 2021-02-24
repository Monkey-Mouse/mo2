<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <editor-menu-bar :editor="editor" v-slot="{ commands, isActive }">
          <div class="menubar">
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
        </editor-menu-bar>
        <v-progress-circular
          v-if="isuploading"
          class="offset-10 offset-lg-7"
          style="position: fixed; z-index: 999"
          indeterminate
          color="amber"
        ></v-progress-circular>
        <editor-menu-bubble
          :editor="editor"
          :keep-in-bounds="true"
          v-slot="{ commands, isActive, menu, getMarkAttrs }"
          class="accent"
        >
          <div
            class="menububble"
            :class="{ 'is-active': menu.isActive }"
            :style="`left: ${menu.left}px; bottom: ${menu.bottom}px;`"
          >
            <div v-if="!linkMenuIsActive">
              <button
                class="menububble__button"
                :class="{ 'is-active': isActive.bold() }"
                @click="commands.bold"
              >
                <v-icon>mdi-format-bold</v-icon>
              </button>

              <button
                class="menububble__button"
                :class="{ 'is-active': isActive.italic() }"
                @click="commands.italic"
              >
                <v-icon>mdi-format-italic</v-icon>
              </button>
              <button
                class="menububble__button"
                :class="{ 'is-active': isActive.underline() }"
                @click="commands.underline"
              >
                <v-icon>mdi-format-underline</v-icon>
              </button>
              <button
                class="menububble__button"
                :class="{ 'is-active': isActive.strike() }"
                @click="commands.strike"
              >
                <v-icon>mdi-format-strikethrough-variant</v-icon>
              </button>

              <button
                class="menububble__button"
                :class="{ 'is-active': isActive.code() }"
                @click="commands.code"
              >
                <v-icon>mdi-code-tags</v-icon>
              </button>
              <button
                v-if="isActive.table()"
                class="menubar__button"
                @click="commands.toggleCellMerge"
              >
                <v-icon>mdi-table-merge-cells</v-icon>/
                <v-icon>mdi-table-split-cell</v-icon>
              </button>
            </div>
            <form
              class="menububble__form"
              v-if="linkMenuIsActive"
              @submit.prevent="setLinkUrl(commands.link, linkUrl)"
            >
              <input
                class="menububble__input"
                type="text"
                v-model="linkUrl"
                placeholder="https://"
                ref="linkInput"
                @keydown.esc="hideLinkMenu"
              />
              <button
                class="menububble__button"
                @click="setLinkUrl(commands.link, null)"
                type="button"
              >
                <v-icon>mdi-close</v-icon>
              </button>
            </form>

            <template v-else>
              <button
                class="menububble__button"
                @click="showLinkMenu(getMarkAttrs('link'))"
                :class="{ 'is-active': isActive.link() }"
              >
                <span>{{ isActive.link() ? "Update Link" : "Add Link" }}</span>
                <v-icon>mdi-link-variant-plus</v-icon>
              </button>
            </template>
          </div>
        </editor-menu-bubble>
        <editor-floating-menu
          :editor="editor"
          v-slot="{ commands, isActive, menu }"
        >
          <div
            class="editor__floating-menu"
            :class="{ 'is-active': menu.isActive }"
            :style="`top: ${menu.top}px`"
          >
            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.heading({ level: 1 }) }"
              @click="commands.heading({ level: 1 })"
            >
              <v-icon>mdi-format-header-1</v-icon>
            </button>

            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.heading({ level: 2 }) }"
              @click="commands.heading({ level: 2 })"
            >
              <v-icon>mdi-format-header-2</v-icon>
            </button>

            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.heading({ level: 3 }) }"
              @click="commands.heading({ level: 3 })"
            >
              <v-icon>mdi-format-header-3</v-icon>
            </button>

            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.bullet_list() }"
              @click="commands.bullet_list"
            >
              <v-icon>mdi-format-list-bulleted</v-icon>
            </button>

            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.ordered_list() }"
              @click="commands.ordered_list"
            >
              <v-icon>mdi-format-list-numbered</v-icon>
            </button>
            <button class="menubar__button" @click="commands.todo_list">
              <v-icon>mdi-format-list-checks</v-icon>
            </button>
            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.blockquote() }"
              @click="commands.blockquote"
            >
              <v-icon>mdi-comment-quote</v-icon>
            </button>

            <button
              class="menubar__button"
              :class="{ 'is-active': isActive.code_block() }"
              @click="commands.code_block"
            >
              <v-icon>mdi-code-braces</v-icon>
            </button>
            <button class="menubar__button" @click="commands.horizontal_rule">
              <v-icon>mdi-minus</v-icon>
            </button>
            <button
              class="menubar__button"
              @click="
                commands.createTable({
                  rowsCount: 3,
                  colsCount: 3,
                  withHeaderRow: false,
                })
              "
            >
              <v-icon>mdi-table-large</v-icon>
            </button>
          </div>
        </editor-floating-menu>
        <div class="mo2content" spellcheck="false">
          <editor-content :editor="editor" />
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
  EditorFloatingMenu,
  EditorMenuBubble,
  EditorMenuBar,
} from "tiptap";
import {
  Blockquote,
  CodeBlock,
  HardBreak,
  Heading,
  OrderedList,
  BulletList,
  ListItem,
  TodoItem,
  TodoList,
  Bold,
  Code,
  Italic,
  Link,
  Strike,
  Underline,
  History,
  CodeBlockHighlight,
  Placeholder,
  TrailingNode,
  HorizontalRule,
  Table,
  TableHeader,
  TableCell,
  TableRow,
  Image,
} from "tiptap-extensions";
import Title from "./title";
import DOC from "./doc";
let that: MO2Editor | any = {};
//#region hljs
import onec from "highlight.js/lib/languages/1c";
import abnf from "highlight.js/lib/languages/abnf";
import accesslog from "highlight.js/lib/languages/accesslog";
import actionscript from "highlight.js/lib/languages/actionscript";
import ada from "highlight.js/lib/languages/ada";
import angelscript from "highlight.js/lib/languages/angelscript";
import apache from "highlight.js/lib/languages/apache";
import applescript from "highlight.js/lib/languages/applescript";
import arcade from "highlight.js/lib/languages/arcade";
import arduino from "highlight.js/lib/languages/arduino";
import armasm from "highlight.js/lib/languages/armasm";
import asciidoc from "highlight.js/lib/languages/asciidoc";
import aspectj from "highlight.js/lib/languages/aspectj";
import autohotkey from "highlight.js/lib/languages/autohotkey";
import autoit from "highlight.js/lib/languages/autoit";
import avrasm from "highlight.js/lib/languages/avrasm";
import awk from "highlight.js/lib/languages/awk";
import axapta from "highlight.js/lib/languages/axapta";
import bash from "highlight.js/lib/languages/bash";
import basic from "highlight.js/lib/languages/basic";
import bnf from "highlight.js/lib/languages/bnf";
import brainfuck from "highlight.js/lib/languages/brainfuck";
import c_like from "highlight.js/lib/languages/c-like";
import c from "highlight.js/lib/languages/c";
import cal from "highlight.js/lib/languages/cal";
import capnproto from "highlight.js/lib/languages/capnproto";
import ceylon from "highlight.js/lib/languages/ceylon";
import clean from "highlight.js/lib/languages/clean";
import clojure_repl from "highlight.js/lib/languages/clojure-repl";
import clojure from "highlight.js/lib/languages/clojure";
import cmake from "highlight.js/lib/languages/cmake";
import coffeescript from "highlight.js/lib/languages/coffeescript";
import coq from "highlight.js/lib/languages/coq";
import cos from "highlight.js/lib/languages/cos";
import cpp from "highlight.js/lib/languages/cpp";
import crmsh from "highlight.js/lib/languages/crmsh";
import crystal from "highlight.js/lib/languages/crystal";
import csharp from "highlight.js/lib/languages/csharp";
import csp from "highlight.js/lib/languages/csp";
import css from "highlight.js/lib/languages/css";
import d from "highlight.js/lib/languages/d";
import dart from "highlight.js/lib/languages/dart";
import delphi from "highlight.js/lib/languages/delphi";
import diff from "highlight.js/lib/languages/diff";
import django from "highlight.js/lib/languages/django";
import dns from "highlight.js/lib/languages/dns";
import dockerfile from "highlight.js/lib/languages/dockerfile";
import dos from "highlight.js/lib/languages/dos";
import dsconfig from "highlight.js/lib/languages/dsconfig";
import dts from "highlight.js/lib/languages/dts";
import dust from "highlight.js/lib/languages/dust";
import ebnf from "highlight.js/lib/languages/ebnf";
import elixir from "highlight.js/lib/languages/elixir";
import elm from "highlight.js/lib/languages/elm";
import erb from "highlight.js/lib/languages/erb";
import erlang_repl from "highlight.js/lib/languages/erlang-repl";
import erlang from "highlight.js/lib/languages/erlang";
import excel from "highlight.js/lib/languages/excel";
import fix from "highlight.js/lib/languages/fix";
import flix from "highlight.js/lib/languages/flix";
import fortran from "highlight.js/lib/languages/fortran";
import fsharp from "highlight.js/lib/languages/fsharp";
import gams from "highlight.js/lib/languages/gams";
import gauss from "highlight.js/lib/languages/gauss";
import gcode from "highlight.js/lib/languages/gcode";
import gherkin from "highlight.js/lib/languages/gherkin";
import glsl from "highlight.js/lib/languages/glsl";
import gml from "highlight.js/lib/languages/gml";
import go from "highlight.js/lib/languages/go";
import golo from "highlight.js/lib/languages/golo";
import gradle from "highlight.js/lib/languages/gradle";
import groovy from "highlight.js/lib/languages/groovy";
import haml from "highlight.js/lib/languages/haml";
import handlebars from "highlight.js/lib/languages/handlebars";
import haskell from "highlight.js/lib/languages/haskell";
import haxe from "highlight.js/lib/languages/haxe";
import hsp from "highlight.js/lib/languages/hsp";
import htmlbars from "highlight.js/lib/languages/htmlbars";
import http from "highlight.js/lib/languages/http";
import hy from "highlight.js/lib/languages/hy";
import inform7 from "highlight.js/lib/languages/inform7";
import ini from "highlight.js/lib/languages/ini";
import irpf90 from "highlight.js/lib/languages/irpf90";
import isbl from "highlight.js/lib/languages/isbl";
import java from "highlight.js/lib/languages/java";
import javascript from "highlight.js/lib/languages/javascript";
import jboss_cli from "highlight.js/lib/languages/jboss-cli";
import json from "highlight.js/lib/languages/json";
import julia_repl from "highlight.js/lib/languages/julia-repl";
import julia from "highlight.js/lib/languages/julia";
import kotlin from "highlight.js/lib/languages/kotlin";
import lasso from "highlight.js/lib/languages/lasso";
import latex from "highlight.js/lib/languages/latex";
import ldif from "highlight.js/lib/languages/ldif";
import leaf from "highlight.js/lib/languages/leaf";
import less from "highlight.js/lib/languages/less";
import lisp from "highlight.js/lib/languages/lisp";
import livecodeserver from "highlight.js/lib/languages/livecodeserver";
import livescript from "highlight.js/lib/languages/livescript";
import llvm from "highlight.js/lib/languages/llvm";
import lsl from "highlight.js/lib/languages/lsl";
import lua from "highlight.js/lib/languages/lua";
import makefile from "highlight.js/lib/languages/makefile";
import markdown from "highlight.js/lib/languages/markdown";
import mathematica from "highlight.js/lib/languages/mathematica";
import matlab from "highlight.js/lib/languages/matlab";
import maxima from "highlight.js/lib/languages/maxima";
import mel from "highlight.js/lib/languages/mel";
import mercury from "highlight.js/lib/languages/mercury";
import mipsasm from "highlight.js/lib/languages/mipsasm";
import mizar from "highlight.js/lib/languages/mizar";
import mojolicious from "highlight.js/lib/languages/mojolicious";
import monkey from "highlight.js/lib/languages/monkey";
import moonscript from "highlight.js/lib/languages/moonscript";
import noneql from "highlight.js/lib/languages/n1ql";
import nginx from "highlight.js/lib/languages/nginx";
import nim from "highlight.js/lib/languages/nim";
import nix from "highlight.js/lib/languages/nix";
import node_repl from "highlight.js/lib/languages/node-repl";
import nsis from "highlight.js/lib/languages/nsis";
import objectivec from "highlight.js/lib/languages/objectivec";
import ocaml from "highlight.js/lib/languages/ocaml";
import openscad from "highlight.js/lib/languages/openscad";
import oxygene from "highlight.js/lib/languages/oxygene";
import parser3 from "highlight.js/lib/languages/parser3";
import perl from "highlight.js/lib/languages/perl";
import pf from "highlight.js/lib/languages/pf";
import pgsql from "highlight.js/lib/languages/pgsql";
import php_template from "highlight.js/lib/languages/php-template";
import php from "highlight.js/lib/languages/php";
import plaintext from "highlight.js/lib/languages/plaintext";
import pony from "highlight.js/lib/languages/pony";
import powershell from "highlight.js/lib/languages/powershell";
import processing from "highlight.js/lib/languages/processing";
import profile from "highlight.js/lib/languages/profile";
import prolog from "highlight.js/lib/languages/prolog";
import properties from "highlight.js/lib/languages/properties";
import protobuf from "highlight.js/lib/languages/protobuf";
import puppet from "highlight.js/lib/languages/puppet";
import purebasic from "highlight.js/lib/languages/purebasic";
import python_repl from "highlight.js/lib/languages/python-repl";
import python from "highlight.js/lib/languages/python";
import q from "highlight.js/lib/languages/q";
import qml from "highlight.js/lib/languages/qml";
import r from "highlight.js/lib/languages/r";
import reasonml from "highlight.js/lib/languages/reasonml";
import rib from "highlight.js/lib/languages/rib";
import roboconf from "highlight.js/lib/languages/roboconf";
import routeros from "highlight.js/lib/languages/routeros";
import rsl from "highlight.js/lib/languages/rsl";
import ruby from "highlight.js/lib/languages/ruby";
import ruleslanguage from "highlight.js/lib/languages/ruleslanguage";
import rust from "highlight.js/lib/languages/rust";
import sas from "highlight.js/lib/languages/sas";
import scala from "highlight.js/lib/languages/scala";
import scheme from "highlight.js/lib/languages/scheme";
import scilab from "highlight.js/lib/languages/scilab";
import scss from "highlight.js/lib/languages/scss";
import shell from "highlight.js/lib/languages/shell";
import smali from "highlight.js/lib/languages/smali";
import smalltalk from "highlight.js/lib/languages/smalltalk";
import sml from "highlight.js/lib/languages/sml";
import sqf from "highlight.js/lib/languages/sqf";
import sql from "highlight.js/lib/languages/sql";
import sql_more from "highlight.js/lib/languages/sql_more";
import stan from "highlight.js/lib/languages/stan";
import stata from "highlight.js/lib/languages/stata";
import step2one from "highlight.js/lib/languages/step21";
import stylus from "highlight.js/lib/languages/stylus";
import subunit from "highlight.js/lib/languages/subunit";
import swift from "highlight.js/lib/languages/swift";
import taggerscript from "highlight.js/lib/languages/taggerscript";
import tap from "highlight.js/lib/languages/tap";
import tcl from "highlight.js/lib/languages/tcl";
import thrift from "highlight.js/lib/languages/thrift";
import tp from "highlight.js/lib/languages/tp";
import twig from "highlight.js/lib/languages/twig";
import typescript from "highlight.js/lib/languages/typescript";
import vala from "highlight.js/lib/languages/vala";
import vbnet from "highlight.js/lib/languages/vbnet";
import vbscript_html from "highlight.js/lib/languages/vbscript-html";
import vbscript from "highlight.js/lib/languages/vbscript";
import verilog from "highlight.js/lib/languages/verilog";
import vhdl from "highlight.js/lib/languages/vhdl";
import vim from "highlight.js/lib/languages/vim";
import x86asm from "highlight.js/lib/languages/x86asm";
import xl from "highlight.js/lib/languages/xl";
import xml from "highlight.js/lib/languages/xml";
import xquery from "highlight.js/lib/languages/xquery";
import yaml from "highlight.js/lib/languages/yaml";
import zephir from "highlight.js/lib/languages/zephir";
import { Prop, Watch } from "vue-property-decorator";
//#endregion
@Component({
  components: {
    EditorContent,
    EditorFloatingMenu,
    EditorMenuBubble,
    EditorMenuBar,
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
  editor: Editor = new Editor({
    extensions: [
      new Blockquote(),
      new CodeBlock(),
      new HardBreak(),
      new Heading({ levels: [1, 2, 3] }),
      new BulletList(),
      new OrderedList(),
      new ListItem(),
      new TodoItem({
        nested: true,
      }),
      new TodoList(),
      new Bold(),
      new Code(),
      new Italic(),
      new Link(),
      new Strike(),
      new Underline(),
      new History(),
      new Placeholder({
        emptyEditorClass: "is-editor-empty",
        emptyNodeClass: "is-empty",
        showOnlyWhenEditable: true,
        showOnlyCurrent: false,
        emptyNodeText: (node) => {
          if (node.type.name === "title") {
            return "Your title";
          }
          return "Your awesome content";
        },
      }),
      new HorizontalRule(),
      new TrailingNode({
        node: "paragraph",
        notAfter: ["paragraph"],
      }),
      new DOC(),
      new Title(),
      new Image(),
      new CodeBlockHighlight({
        languages: {
          onec,
          abnf,
          accesslog,
          actionscript,
          ada,
          angelscript,
          apache,
          applescript,
          arcade,
          arduino,
          armasm,
          asciidoc,
          aspectj,
          autohotkey,
          autoit,
          avrasm,
          awk,
          axapta,
          bash,
          basic,
          bnf,
          brainfuck,
          c_like,
          c,
          cal,
          capnproto,
          ceylon,
          clean,
          clojure_repl,
          clojure,
          cmake,
          coffeescript,
          coq,
          cos,
          cpp,
          crmsh,
          crystal,
          csharp,
          csp,
          css,
          d,
          dart,
          delphi,
          diff,
          django,
          dns,
          dockerfile,
          dos,
          dsconfig,
          dts,
          dust,
          ebnf,
          elixir,
          elm,
          erb,
          erlang_repl,
          erlang,
          excel,
          fix,
          flix,
          fortran,
          fsharp,
          gams,
          gauss,
          gcode,
          gherkin,
          glsl,
          gml,
          go,
          golo,
          gradle,
          groovy,
          haml,
          handlebars,
          haskell,
          haxe,
          hsp,
          htmlbars,
          http,
          hy,
          inform7,
          ini,
          irpf90,
          isbl,
          java,
          javascript,
          jboss_cli,
          json,
          julia_repl,
          julia,
          kotlin,
          lasso,
          latex,
          ldif,
          leaf,
          less,
          lisp,
          livecodeserver,
          livescript,
          llvm,
          lsl,
          lua,
          makefile,
          markdown,
          mathematica,
          matlab,
          maxima,
          mel,
          mercury,
          mipsasm,
          mizar,
          mojolicious,
          monkey,
          moonscript,
          noneql,
          nginx,
          nim,
          nix,
          node_repl,
          nsis,
          objectivec,
          ocaml,
          openscad,
          oxygene,
          parser3,
          perl,
          pf,
          pgsql,
          php_template,
          php,
          plaintext,
          pony,
          powershell,
          processing,
          profile,
          prolog,
          properties,
          protobuf,
          puppet,
          purebasic,
          python_repl,
          python,
          q,
          qml,
          r,
          reasonml,
          rib,
          roboconf,
          routeros,
          rsl,
          ruby,
          ruleslanguage,
          rust,
          sas,
          scala,
          scheme,
          scilab,
          scss,
          shell,
          smali,
          smalltalk,
          sml,
          sqf,
          sql,
          sql_more,
          stan,
          stata,
          step2one,
          stylus,
          subunit,
          swift,
          taggerscript,
          tap,
          tcl,
          thrift,
          tp,
          twig,
          typescript,
          vala,
          vbnet,
          vbscript_html,
          vbscript,
          verilog,
          vhdl,
          vim,
          x86asm,
          xl,
          xml,
          xquery,
          yaml,
          zephir,
        },
      }),
      new Table({
        resizable: true,
      }),
      new TableHeader(),
      new TableCell(),
      new TableRow(),
    ],
    content: `
    `,
    onUpdate() {
      (that as MO2Editor).update = true;
    },
    onPaste(editorview, event, slice) {
      var items = (event.clipboardData || event.originalEvent.clipboardData)
        .items;
      var files = [];
      for (let index = 0; index < items.length; index++) {
        var item = items[index];
        if (item.kind === "file") {
          var blob = item.getAsFile();
          files = files.concat(blob);
          // var reader = new FileReader();
          // reader.onload = function(event){
          //   console.log(event.target.result)}; // data url!
          // reader.readAsDataURL(blob);
        }
      }
      if (files.length === 0) {
        return;
      }
      (event as Event).preventDefault();
      that.isuploading = true;
      that
        .uploadImgs(files, this.commands.image)
        .then(() => (that.isuploading = false))
        .catch(() => {
          that.isuploading = false;
        });
    },
  });
  mounted() {
    that = this;
    if (this.content) {
      this.editor.setContent(this.content);
    }
    this.$emit("loaded", this);
    setInterval(() => {
      if ((that as MO2Editor).update) {
        (that as MO2Editor).update = false;
        (that as MO2Editor).$emit("autosave");
      }
    }, 5000);
  }

  @Watch("content")
  contentSet() {
    if (this.content) {
      this.editor.setContent(this.content);
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