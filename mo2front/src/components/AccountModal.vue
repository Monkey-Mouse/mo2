<template>
  <v-dialog
    :value="enable"
    @click:outside="close"
    max-width="600px"
    autocomplete="off"
  >
    <v-card>
      <v-container>
        <v-row>
          <v-col cols="12">
            <v-card-title>
              <v-tabs align-with-title v-model="tabkey">
                <v-tab :key="1">登录</v-tab>
                <v-tab :key="2">注册</v-tab>
              </v-tabs>
            </v-card-title>
          </v-col>
        </v-row>
        <v-tabs-items v-model="tabkey">
          <v-tab-item :key="1">
            <v-card-text>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    label="Email"
                    v-model="email"
                    :rules="validateEmail()"
                  >
                    <v-icon slot="append" color="gray"> mdi-email </v-icon>
                  </v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    label="Password"
                    v-model="password"
                    :type="showPasswd ? 'text' : 'password'"
                    :append-icon="showPasswd ? 'mdi-eye-off' : 'mdi-eye'"
                    @click:append="showPasswd = !showPasswd"
                    hint="长度最小为8"
                    :rules="validatePasswd()"
                  >
                  </v-text-field>
                </v-col>
              </v-row>
            </v-card-text>
            <v-card-actions>
              <v-switch label="记住我"></v-switch>
              <v-spacer></v-spacer>
              <v-btn outlined text>{{ reg ? "注册&" : "" }}登录</v-btn>
              <v-btn @click="close" color="red">取消</v-btn>
            </v-card-actions>
          </v-tab-item>
          <v-tab-item :key="2">
            <v-card-text>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    label="Name"
                    v-model="name"
                    :rules="validateName()"
                  >
                    <v-icon slot="append" color="gray"> mdi-account </v-icon>
                  </v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    label="Email"
                    v-model="email"
                    :rules="validateEmail()"
                  >
                    <v-icon slot="append" color="gray"> mdi-email </v-icon>
                  </v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    label="Password"
                    v-model="password"
                    :type="showPasswd ? 'text' : 'password'"
                    :append-icon="showPasswd ? 'mdi-eye-off' : 'mdi-eye'"
                    @click:append="showPasswd = !showPasswd"
                    hint="长度最小为8"
                    :rules="validatePasswd()"
                  >
                  </v-text-field>
                </v-col>
              </v-row>
            </v-card-text>
            <v-card-actions>
              <v-switch label="记住我"></v-switch>
              <v-spacer></v-spacer>
              <v-btn outlined text>注册&登录</v-btn>
              <v-btn @click="close" color="red">取消</v-btn>
            </v-card-actions>
          </v-tab-item>
        </v-tabs-items>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import {
  required,
  minLength,
  maxLength,
  email,
} from "vuelidate/lib/validators";

@Component({})
export default class AccountModal extends Vue {
  @Prop()
  enable!: boolean;
  email: string = "";
  name: string = "";
  password: string = "";
  tabkey = 1;
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
  reg = false;
  showPasswd: boolean = false;
  created() {
    this.email = "";
    this.password = "";
  }
  close() {
    this.$emit("update:enable", false);
  }
  validateName() {
    this.$v.name.$touch();

    return [() => this.$v.name.required || "名字不可为空"];
  }

  validatePasswd() {
    this.$v.password.$touch();

    return [
      () => this.$v.password.required || "密码不可为空",
      () => this.$v.password.min || "密码长度不小于8",
    ];
  }

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
