<template>
  <v-dialog :value="show" @click:outside="close" max-width="600px">
    <v-card>
      <v-container>
        <v-row>
          <v-col cols="12">
            <v-card-title> {{ title }} </v-card-title>
          </v-col>
        </v-row>
        <v-card-text>
          <input-list :validator="validator" :inputProps="inputProps" />
        </v-card-text>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import InputList from "../components/InputList.vue";
import Cropper from "../components/ImageCropper.vue";
import { InputProp } from "@/models";
import { Prop } from "vue-property-decorator";
@Component({
  components: {
    InputList,
    Cropper,
  },
})
export default class About extends Vue {
  @Prop()
  title!: string;
  @Prop()
  validator!: any;
  @Prop()
  inputProps!: { [name: string]: InputProp };
  @Prop()
  show!: boolean;
  close() {
    this.$emit("update:show", false);
  }
  // validator = {
  //   password: {
  //     required: required,
  //     min: minLength(8),
  //   },
  //   email: {
  //     required: required,
  //     email: email,
  //   },
  //   name: {
  //     required: required,
  //   },
  // };
  // inputProps: { [name: string]: InputProp } = {
  //   email: {
  //     errorMsg: {
  //       required: "email不可为空",
  //       email: "请输入合法的email",
  //     },
  //     label: "Email",
  //     default: "",
  //     icon: "mdi-email",
  //     col: 6,
  //     type: "email",
  //   },
  //   password: {
  //     errorMsg: {
  //       required: "password不可为空",
  //       min: "password不可短于8",
  //     },
  //     label: "Password",
  //     default: "",
  //     icon: "mdi-key",
  //     col: 6,
  //     type: "password",
  //     iconClick: (prop) => {
  //       if (prop.type === "text") {
  //         prop.type = "password";
  //       } else prop.type = "text";
  //     },
  //   },
  // };
}
</script>
