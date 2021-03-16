<template>
  <v-container>
    <v-row>
      <v-col cols="12" lg="8" offset-lg="2">
        <input-list :inputProps="inputProps" :validator="validator" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { InputProp } from "@/models";
import { required } from "vuelidate/lib/validators";
import Vue from "vue";
import Component from "vue-class-component";
import InputList from "../components/InputList.vue";
@Component({
  components: { InputList },
})
export default class Settings extends Vue {
  validator = {
    name: {
      required: required,
    },
  };
  inputProps: { [name: string]: InputProp } = {
    name: {
      errorMsg: {
        required: "用户名不可为空",
      },
      label: "Name",
      default: "",
      icon: "mdi-account",
      col: 6,
      type: "text",
    },
    darkMode: {
      label: "prefer dark",
      default: false,
      icon: "mdi-theme-light-dark",
      col: 6,
      type: "switch",
      errorMsg: {},
    },
  };
}
</script>