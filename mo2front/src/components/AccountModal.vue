<template>
  <div class="modal is-active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <div class="modal-card-body">
        <v-text-field label="Email" v-model="email" :rules="validateEmail()">
          <v-icon slot="append" color="gray"> mdi-email </v-icon>
        </v-text-field>
        <v-text-field label="Password" v-model="email" :rules="validateEmail()">
          <v-icon slot="append" color="gray"> mdi-email </v-icon>
        </v-text-field>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {
  required,
  minLength,
  maxLength,
  email,
} from "vuelidate/lib/validators";

@Component({})
export default class AccountModal extends Vue {
  email: string = "";
  validator = {
    email: {
      required: required,
      email: email,
    },
  };
  validateEmail() {
    this.$v.email.$touch();

    return [
      () => this.$v.email.required || "email不可为空",
      () => this.$v.email.email || "请输入一个正确的email值",
    ];
  }
  validations() {
    return this.validator;
  }
}
</script>
