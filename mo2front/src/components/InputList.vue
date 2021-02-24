<template>
  <v-card-text>
    <v-row>
      <v-col v-for="(value, key) in inputProps" :key="key" :cols="value['col']">
        <v-textarea
          v-if="value['type'] === 'textarea'"
          outlined
          auto-grow
          :label="value.label"
          v-model="model[key]"
          :rules="buildValidationRoles(key)"
          :type="value['type']"
          :append-icon="value.icon"
        />
        <img-selector
          v-else-if="value['type'] === 'imgselector'"
          :imgs="imgVals[key]"
          :title="value.label"
          :appendIcon="value.icon"
          :uploadImgs="uploadImgs"
          @imgselect="(img) => (model[key] = img)"
        />
        <v-text-field
          v-else
          :label="value.label"
          v-model="model[key]"
          :rules="buildValidationRoles(key)"
          :type="value['type']"
        >
          <v-icon
            class="clickable"
            v-if="value.icon && value['iconClick']"
            slot="append"
            color="gray"
            @click="value['iconClick'](value)"
          >
            {{ value.icon }}
          </v-icon>
          <v-icon v-else-if="value.icon" slot="append" color="gray">
            {{ value.icon }}
          </v-icon>
        </v-text-field>
      </v-col>
    </v-row>
  </v-card-text>
</template>

<script lang="ts">
import { InputProp } from "@/models";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import ImgSelector from "./ImgSelector.vue";

@Component({
  components: {
    ImgSelector,
  },
})
export default class InputList extends Vue {
  model: any = {};
  @Prop()
  validator!: any;
  @Prop()
  inputProps!: { [name: string]: InputProp };
  @Prop()
  anyError: boolean;
  @Prop()
  uploadImgs: (
    blobs: File[],
    callback: (imgprop: { src: string }) => void
  ) => Promise<void>;
  imgVals: any = {};
  constructor() {
    super();
    for (let key in this.inputProps) {
      if (this.inputProps[key].type === "imgselector") {
        this.imgVals[key] = this.inputProps[key].default;
      } else this.model[key] = this.inputProps[key].default;
    }
  }
  mounted() {
    this.$emit("update:anyError", this.$v.$anyError);
  }
  buildValidationRoles(prop: string) {
    const errmsgs = this.inputProps[prop].errorMsg;
    const rules = [];
    const vuelidator = this.$v["model"];
    for (let key in errmsgs) {
      rules.push(() => {
        vuelidator[prop].$touch();
        return vuelidator[prop][key] || errmsgs[key];
      });
    }
    return rules;
  }
  @Watch("$v.$anyError")
  errorChange() {
    this.$emit("update:anyError", this.$v.$anyError);
  }
  setModel(model: any) {
    for (const key in this.model) {
      if (model[key]) {
        this.model[key] = model[key];
      }
    }
    for (const key in this.imgVals) {
      if (model[key]) {
        this.imgVals[key] = model[key];
        for (let index = 0; index < model[key].length; index++) {
          const element = model[key][index];
          if (element.active) {
            this.model[key] = element.src;
          }
        }
      }
    }
  }

  validations() {
    return { model: this.validator };
  }
  public get Model() {
    return this.model;
  }
}
</script>

