<template>
  <v-dialog :value="show" @click:outside="close" max-width="600px">
    <v-card
      @input="confirmerr = ''"
      :disabled="processing"
      :loading="processing"
    >
      <v-container>
        <v-row>
          <v-col cols="12">
            <v-card-title> {{ title }} </v-card-title>
          </v-col>
        </v-row>
        <v-card-text>
          <input-list
            :anyError.sync="anyError"
            ref="inputs"
            :validator="validator"
            :inputProps="inputProps"
            :uploadImgs="uploadImgs"
          />
          <v-row v-if="confirmerr !== ''">
            <v-alert dense outlined type="error" class="col-12">{{
              confirmerr
            }}</v-alert></v-row
          >
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="anyError"
            color="green"
            outlined
            text
            @click="confirmClick"
            >{{ confirmText }}</v-btn
          >
          <v-btn @click="close" color="red">取消</v-btn>
        </v-card-actions>
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
  confirm!: (any) => Promise<{ err: string; pass: boolean }>;
  @Prop()
  title!: string;
  @Prop()
  validator!: any;
  @Prop()
  inputProps!: { [name: string]: InputProp };
  @Prop()
  show!: boolean;
  @Prop()
  confirmText!: string;
  @Prop()
  uploadImgs: (
    blobs: File[],
    callback: (imgprop: { src: string }) => void
  ) => Promise<void>;
  confirmerr = "";
  anyError = true;
  processing = false;
  close() {
    this.$emit("update:show", false);
  }
  async confirmClick() {
    this.processing = true;
    const { err, pass } = await this.confirm(
      (this.$refs["inputs"] as InputList).Model
    );
    this.processing = false;
    if (pass) {
      this.close();
    } else {
      this.confirmerr = err;
    }
  }
  setModel(model: any) {
    setTimeout(() => {
      (this.$refs["inputs"] as InputList).setModel(model);
    }, 100);
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
