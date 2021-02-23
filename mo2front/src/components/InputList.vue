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

@Component({})
export default class InputList extends Vue {
  model: any = {};
  @Prop()
  validator!: any;
  @Prop()
  inputProps!: { [name: string]: InputProp };
  @Prop()
  anyError: boolean;
  constructor() {
    super();
    for (let key in this.inputProps) {
      this.model[key] = this.inputProps[key].default;
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
      this.model[key] = model[key];
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

