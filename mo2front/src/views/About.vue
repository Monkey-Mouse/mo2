<template>
  <div class="about">
    <h1>This is an about page</h1>
    <v-dialog :value="false" max-width="600px">
      <v-card>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-card-title> FUCK </v-card-title>
            </v-col>
          </v-row>
          <v-card-text>
            <input-list :validator="validator" :inputProps="inputProps" />
          </v-card-text>
        </v-container>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import InputList from "../components/InputList.vue";
import Cropper from "../components/ImageCropper.vue";
import {
  required,
  minLength,
  maxLength,
  email,
} from "vuelidate/lib/validators";
import { InputProp } from "@/models";
@Component({
  components: {
    InputList,
    Cropper,
  },
})
export default class About extends Vue {
  validator = {
    password: {
      required: required,
      min: minLength(8),
    },
    email: {
      required: required,
      email: email,
    },
    name: {
      required: required,
    },
  };
  inputProps: { [name: string]: InputProp } = {
    email: {
      errorMsg: {
        required: "email不可为空",
        email: "请输入合法的email",
      },
      label: "Email",
      default: "",
      icon: "mdi-email",
      col: 6,
      type: "email",
    },
    password: {
      errorMsg: {
        required: "password不可为空",
        min: "password不可短于8",
      },
      label: "Password",
      default: "",
      icon: "mdi-key",
      col: 6,
      type: "password",
      iconClick: (prop) => {
        if (prop.type === "text") {
          prop.type = "password";
        } else prop.type = "text";
      },
    },
  };
}
</script>
