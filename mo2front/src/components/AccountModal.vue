<template>
  <div>
    <v-snackbar v-model="snackbar" :timeout="5000">
      {{ "确认Email即将发送到你的邮箱，请点击邮箱中的确认按钮后再继续！" }}

      <template v-slot:action="{ attrs }">
        <v-btn color="blue" text v-bind="attrs" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
    <v-snackbar v-model="tipbar" :timeout="5000">
      {{ "你可以在确认邮箱之后再次点击右上角的登录继续" }}

      <template v-slot:action="{ attrs }">
        <v-btn color="blue" text v-bind="attrs" @click="tipbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
    <v-dialog
      :value="enable"
      @click:outside="close"
      max-width="600px"
      autocomplete="off"
    >
      <v-card :loading="processing" :disabled="processing">
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-card-title>
                <v-tabs align-with-title v-model="tabkey">
                  <v-tab :key="1">登录</v-tab>
                  <v-tab :key="2">注册</v-tab>
                  <v-tab :key="2">OAuth</v-tab>
                </v-tabs>
              </v-card-title>
            </v-col>
          </v-row>
          <v-tabs-items v-model="tabkey">
            <v-tab-item @focus="regerror = ''" :key="1">
              <v-card-text>
                <v-row>
                  <v-col cols="12">
                    <v-text-field
                      label="Email or UserName"
                      v-model="emailOrName"
                      :rules="validateNameOrEmail()"
                    >
                      <v-icon slot="append" color="gray"> mdi-account </v-icon>
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
                <v-row v-if="loginerr !== ''">
                  <v-alert dense outlined type="error" class="col-12">{{
                    loginerr
                  }}</v-alert></v-row
                >
              </v-card-text>
              <v-card-actions>
                <!-- <v-switch label="记住我"></v-switch> -->
                <v-spacer></v-spacer>
                <v-btn
                  outlined
                  text
                  :disabled="
                    this.$v.password.$anyError || this.$v.emailOrName.$anyError
                  "
                  @click="login"
                  >登录</v-btn
                >
                <v-btn @click="close" color="error">取消</v-btn>
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
                <v-row v-if="regerror !== ''">
                  <v-alert dense outlined type="error" class="col-12">{{
                    regerror
                  }}</v-alert></v-row
                >
              </v-card-text>
              <v-card-actions>
                <!-- <v-switch label="记住我"></v-switch> -->
                <v-spacer></v-spacer>
                <v-btn
                  v-if="emailSent"
                  :disabled="this.$v.$anyError"
                  outlined
                  text
                  @click="login"
                  >确认</v-btn
                >
                <!-- <v-btn
                  v-if="emailSent"
                  :disabled="this.$v.$anyError"
                  outlined
                  text
                  @click="backToMain"
                  >先继续浏览</v-btn
                > -->
                <v-btn :disabled="regDisable" outlined text @click="register">{{
                  emailSent
                    ? "重新发送" + (seconds > 0 ? `${seconds}` : "")
                    : "注册"
                }}</v-btn>

                <v-btn @click="close" color="error">取消</v-btn>
              </v-card-actions>
            </v-tab-item>
            <v-tab-item :key="2" style="min-height: 150px">
              <v-row justify="center">
                <v-col cols="12" align-self="center" class="text-center"
                  ><v-btn @click="github"
                    ><v-icon>mdi-github</v-icon>使用GitHub账号</v-btn
                  ></v-col
                >
              </v-row>
            </v-tab-item>
          </v-tabs-items>
        </v-container>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { User } from "@/models";
import { GetErrorMsg, GithubOauth, LoginAsync, RegisterAsync } from "../utils";
import { AxiosError } from "axios";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import {
  required,
  minLength,
  maxLength,
  email,
} from "vuelidate/lib/validators";
import { VForm } from "vuetify/lib";

@Component({})
export default class AccountModal extends Vue {
  @Prop()
  enable!: boolean;
  @Prop()
  user!: User;
  emailOrName = "";
  regerror: string = "";
  loginerr = "";
  processing = false;
  email: string = "";
  name: string = "";
  password: string = "";
  emailSent = false;
  snackbar = false;
  tipbar = false;
  tabkey = 0;
  seconds = -1;
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
    emailOrName: {
      required: required,
    },
  };
  showPasswd: boolean = false;

  public get regDisable(): boolean {
    return (
      this.$v.password.$anyError ||
      this.$v.email.$anyError ||
      this.$v.name.$anyError ||
      this.seconds > 0
    );
  }

  created() {
    this.email = "";
    this.password = "";
    setInterval(() => {
      this.seconds--;
    }, 1000);
  }
  backToMain() {
    this.snackbar = false;
    this.tipbar = true;
    this.close();
  }
  github() {
    GithubOauth();
  }
  close() {
    this.$emit("update:enable", false);
  }
  login() {
    if (this.tabkey === 1) {
      this.emailOrName = this.email;
    }
    this.$v.$touch();
    if (this.$v.password.$anyError || this.$v.emailOrName.$anyError) return;
    this.processing = true;
    LoginAsync({
      userNameOrEmail: this.emailOrName,
      password: this.password,
    })
      .then((u) => {
        this.$emit("update:user", u);
        this.close();
        this.processing = false;
      })
      .catch((err) => {
        this.processing = false;
        if (this.emailSent) {
          this.regerror = GetErrorMsg(err);
          return;
        }
        this.loginerr = GetErrorMsg(err);
      });
  }
  validateNameOrEmail() {
    this.$v.emailOrName.$touch();

    return [() => this.$v.emailOrName.required || "登录信息不可为空"];
  }
  validateName() {
    this.$v.name.$touch();

    return [() => this.$v.name.required || "名字不可为空"];
  }
  register() {
    this.$v.$touch();
    if (this.$v.$anyError) return;
    this.processing = true;
    RegisterAsync({
      email: this.email,
      userName: this.name,
      password: this.password,
    })
      .then((u) => {
        this.seconds = 30;
        this.tipbar = false;
        this.snackbar = true;
        this.processing = false;
        this.emailSent = true;
        // SendVerifyEmail(u.email).then(() => {
        //   alert("发送成功");
        //   this.processing = false;
        // });
      })
      .catch((err) => {
        this.processing = false;
        this.regerror = GetErrorMsg(err);
      });
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
