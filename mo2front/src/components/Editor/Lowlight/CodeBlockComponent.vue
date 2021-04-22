<template>
  <node-view-wrapper class="code-block">
    <!-- <select contenteditable="false" v-model="selectedLanguage">
      <option :value="null">auto</option>
      <option disabled>—</option>
      <option
        v-for="(language, index) in languages"
        :value="language"
        :key="index"
      >
        {{ language }}
      </option>
    </select> -->
    <v-autocomplete
      dense
      background-color="primary"
      :items="languages"
      contenteditable="false"
      v-model="selectedLanguage"
      item-text="text"
      item-value="value"
    ></v-autocomplete>
    <pre><node-view-content as="code" /></pre>
  </node-view-wrapper>
</template>

<script>
import { NodeViewWrapper, NodeViewContent, nodeViewProps } from "@tiptap/vue-2";
let first = true;
export default {
  components: {
    NodeViewWrapper,
    NodeViewContent,
  },
  props: nodeViewProps,
  data() {
    const arr = this.extension.options.lowlight.listLanguages().map((v) => {
      return { text: v, value: v };
    });
    arr.push({ text: "auto", value: null });
    return {
      languages: arr,
    };
  },
  computed: {
    selectedLanguage: {
      get() {
        let language = this.node.attrs.language;
        if (language != null) {
          language = language.split(" ")[0];
        }

        if (first) {
          this.updateAttributes({ language });
          first = false;
        }

        return language;
      },
      set(language) {
        this.updateAttributes({ language });
      },
    },
  },
};
</script>

<style lang="scss" scoped>
.code-block {
  position: relative;
  .v-select {
    position: absolute;
    // top: 0.5rem;
    right: 0.5rem;
    padding: 0;
    margin: 0;
    width: 100px;
  }
}
</style>