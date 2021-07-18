<template>
  <v-card-text>
    <v-row>
      <v-col
        v-for="(value, key) in inputProps"
        :key="key"
        :lg="value['col']"
        cols="12"
      >
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
        <div v-else-if="value['type'] === 'color'">
          <div class="v-label mb-3">{{ value.label }}</div>
          <v-color-picker
            :label="value.label"
            v-model="model[key]"
            @update:color="value.onChange"
            mode="hexa"
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
          </v-color-picker>
        </div>
        <v-switch
          v-else-if="value['type'] === 'switch'"
          :label="value.label"
          v-model="model[key]"
          @change="value.onChange"
          :messages="value.message"
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
        </v-switch>
        <img-selector
          v-else-if="value['type'] === 'imgselector'"
          :imgs="imgVals[key]"
          :title="value.label"
          :appendIcon="value.icon"
          :uploadImgs="uploadImgs"
          @imgselect="(img) => (model[key] = img)"
        />
        <v-autocomplete
          v-else-if="value['type'] === 'select'"
          v-model="model[key]"
          :items="value['options']"
          deletable-chips
          chips
          :label="value.label"
          :multiple="value.multiple"
          :loading="value.loading"
          :filter="value.filter"
          @update:search-input="
            (s) => {
              if (value.onChange)
                value.onChange({ input: s, val: value, cu: model[key] });
            }
          "
          :search-input="value.input"
        >
          <template v-if="value['showAvatar']" v-slot:selection="data">
            <v-chip
              v-bind="data.attrs"
              :input-value="data.selected"
              close
              @click="data.select"
              @click:close="
                () => {
                  let index = -1;
                  if (data.item.value)
                    index = model[key].indexOf(data.item.value);
                  else index = model[key].indexOf(data.item);
                  if (index >= 0) model[key].splice(index, 1);
                }
              "
            >
              <v-avatar left>
                <v-img
                  v-if="data.item.avatar"
                  :src="data.item.avatar + '~thumb'"
                ></v-img>
                <span v-else class="white--text headline">{{
                  initials(data.item.text)
                }}</span>
              </v-avatar>
              {{ data.item.text }}
            </v-chip>
          </template>
          <template v-if="value['showAvatar']" v-slot:item="data">
            <template>
              <v-list-item-avatar @click="value.input = ''">
                <img :src="data.item.avatar" />
              </v-list-item-avatar>
              <v-list-item-content @click="value.input = ''">
                <v-list-item-title v-text="data.item.text"></v-list-item-title>
              </v-list-item-content>
            </template>
          </template>
        </v-autocomplete>
        <v-combobox
          v-else-if="value['type'] === 'combo'"
          v-model="model[key]"
          :items="value['options']"
          deletable-chips
          chips
          :label="value.label"
          :multiple="value.multiple"
          :messages="value.message"
        />
        <v-file-input
          v-else-if="value['type'] === 'file'"
          :label="value.label"
          v-model="model[key]"
          :rules="buildValidationRoles(key)"
          :accept="value['accept']"
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
        </v-file-input>
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
import { InputProp } from "../models";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import ImgSelector from "./ImgSelector.vue";
import { GetInitials } from "@/utils";

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
  initials(s: string) {
    return GetInitials(s);
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
      if (model[key] || model[key] === false) {
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

