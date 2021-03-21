<template>
  <div>
    <MO2Dialog
      :validator="validator"
      :inputProps="inputProps"
      :show.sync="edit"
      title="更改信息"
      confirmText="确认"
      :confirm="confirm"
    />
    <cropper
      :show.sync="showCropper"
      :img="avatar"
      title="裁剪你的头像"
      @confirm="confirmAvatar"
      :loading="avatarProcessing"
      :confirmerr="avErr"
      @imgLoad="imgLoad"
    />
    <input
      @change="imgChange"
      ref="f"
      type="file"
      accept="image/*"
      style="display: none"
    />
    <v-parallax
      dark
      src="https://cdn.mo2.leezeeyee.com/material.jpg"
      height="400"
    >
      <v-row align="center" justify="center">
        <v-col class="text-center" cols="12">
          <b @click="changeAvatar" :class="ownPage ? 'clickable' : ''">
            <avatar :size="80" :user="displayUser" />
          </b>
          <!-- <v-img class="v-avatar" :src="displayUser.avatar"></v-img> -->
          <h1 class="display-1 font-weight-thin mb-4">
            {{ displayUser.name }}
            <v-icon v-if="ownPage" @click="edit = true"
              >mdi-account-edit</v-icon
            >
          </h1>
          <h4 class="subheading">{{ displayUser.description }}</h4>
          <h4 class="subtitle-2">
            {{ displayUser.email }}<v-icon color="grey"> mdi-email</v-icon>
          </h4>
        </v-col>
      </v-row>
    </v-parallax>
    <v-container>
      <v-tabs
        v-model="tab"
        background-color=" grey accent-4"
        centered
        dark
        icons-and-text
      >
        <v-tabs-slider></v-tabs-slider>

        <v-tab href="#tab-1"> Articles </v-tab>

        <v-tab v-if="ownPage" href="#tab-2"> DraftBox </v-tab>
        <v-tab href="#tab-3"> Categories </v-tab>
      </v-tabs>

      <v-tabs-items v-model="tab">
        <v-tab-item :value="'tab-1'">
          <v-card flat>
            <blog-time-line-list v-if="!firstloading" :blogs="datalist" />
            <blog-skeleton v-if="loading" :num="pagesize" />
          </v-card>
        </v-tab-item>
        <v-tab-item :value="'tab-2'">
          <v-card flat>
            <blog-time-line-list
              v-if="!draftProps.firstloading"
              :blogs="draftProps.datalist"
              :draft="true"
            />
            <blog-skeleton v-if="draftProps.loading" :num="pagesize" />
          </v-card>
        </v-tab-item>
        <v-tab-item :value="'tab-3'">
          <v-card flat>
            <Category v-if="!firstloading" :own="ownPage" :user="displayUser" />
          </v-card>
        </v-tab-item>
      </v-tabs-items>
    </v-container>
  </div>
</template>

<script lang="ts">
import { BlankUser, BlogBrief, User, InputProp } from "@/models";
import {
  AddMore,
  addQuery,
  AutoLoader,
  ElmReachedButtom,
  GetErrorMsg,
  GetOwnArticles,
  GetUserArticles,
  GetUserData,
  UpdateUserInfo,
  UploadImgToQiniu,
} from "@/utils";
import { required } from "vuelidate/lib/validators";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import Avatar from "../components/UserAvatar.vue";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
import MO2Dialog from "../components/MO2Dialog.vue";
import Cropper from "../components/ImageCropper.vue";
import Category from "../components/Category.vue";
@Component({
  components: {
    BlogTimeLineList,
    Avatar,
    BlogSkeleton,
    MO2Dialog,
    Cropper,
    Category,
  },
})
export default class Account extends Vue implements AutoLoader<BlogBrief> {
  @Prop()
  user!: User;
  displayUser: User = BlankUser;
  uid!: string;
  datalist: BlogBrief[] = [];
  loading = true;
  firstloading = true;
  page = 0;
  pagesize = 5;
  nomore = false;
  tab = "tab-1";
  ownPage = false;
  create = false;
  edit = false;
  avatar: string = "";
  showCropper = false;
  avatarProcessing = false;
  avErr = "";

  draftProps: AutoLoader<BlogBrief> = {
    loading: true,
    firstloading: true,
    page: 0,
    pagesize: 5,
    nomore: false,
    datalist: [],
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    ReachedButtom: () => {},
  };
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
      col: 12,
      type: "text",
    },
  };
  imgLoad(success) {
    if (success === "error") {
      this.avErr = "图片格式错误！";
    }
    this.avatarProcessing = false;
  }
  async confirm({ name }: { name: string }) {
    try {
      this.user.name = name;
      await UpdateUserInfo(this.user);
      // this.displayUser.name = name;
      this.updateUser();
      // this.initPage();
      return { err: "", pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  async confirmAvatar(data: any) {
    this.avatarProcessing = true;
    data.lastModifiedDate = new Date();
    data.name = "avatar.webp";
    let f = data as File;
    if (!this.user.settings) {
      this.user.settings = {};
    }
    try {
      await UploadImgToQiniu(
        [f],
        ({ src }) => (this.user.settings.avatar = src)
      );
    } catch (error) {
      this.avatarProcessing = false;
      this.avErr = "图片上传失败！";
      return;
    }
    try {
      await UpdateUserInfo(this.user);
    } catch (error) {
      this.avatarProcessing = false;
      this.avErr = GetErrorMsg(error);
      return;
    }
    this.showCropper = false;
    // this.displayUser.settings.avatar = this.user.settings.avatar;
    this.updateUser();
    this.avatarProcessing = false;
    this.avErr = "";
  }
  @Watch("showCropper")
  changeCropper() {
    if (!this.showCropper) {
      this.avErr = "";
    }
  }
  updateUser() {
    this.$emit("update:user", this.user);
  }
  changeAvatar() {
    if (!this.ownPage) {
      return;
    }
    (this.$refs.f as HTMLInputElement).click();
  }
  imgChange() {
    this.showCropper = true;
    this.avatarProcessing = true;
    this.avatar = URL.createObjectURL(
      (this.$refs.f as HTMLInputElement).files[0]
    );
    (this.$refs.f as HTMLInputElement).value = "";
  }
  created() {
    this.initPage();
    this.create = true;
  }
  async initPage() {
    if (this.$route.query["tab"]) {
      this.tab = this.$route.query["tab"] as string;
    }
    this.uid = this.$route.params["id"];
    this.inputProps["name"].default = this.user.name;
    if (this.uid === undefined || this.uid === this.user.id) {
      this.uid = this.user.id;
      this.displayUser = this.user;
      this.ownPage = true;
      if (this.$route.path !== "/account") {
        this.$router.replace("/account");
      }
      GetOwnArticles({
        page: this.page++,
        pageSize: this.pagesize,
        draft: false,
      }).then((data) => {
        AddMore(this, data);
        this.firstloading = false;
      });
    } else {
      GetUserData(this.uid)
        .then((u) => {
          this.displayUser = u;
          GetUserArticles({
            page: this.page++,
            pageSize: this.pagesize,
            draft: false,
            id: this.uid,
          }).then((data) => {
            AddMore(this, data);
            this.firstloading = false;
          });
        })
        .catch((err) => GetErrorMsg(err));
    }
  }

  @Watch("$route", { immediate: true, deep: true })
  pageChange() {
    if (
      (!this.$route.params["id"] && this.uid === this.user.id) ||
      this.uid === this.$route.params["id"]
    ) {
      return;
    }
    if (this.create) {
      this.initPage();
    }
  }

  @Watch("user")
  userChange() {
    this.uid = this.user.id;
    this.displayUser = this.user;
    this.inputProps.name.default = this.user.name;
  }
  @Watch("tab")
  loadDraft() {
    addQuery(this, "tab", this.tab);
    if (this.draftProps.firstloading && this.tab === "tab-2") {
      if (this.ownPage) {
        GetOwnArticles({
          page: this.draftProps.page++,
          pageSize: this.draftProps.pagesize,
          draft: true,
        }).then((val) => {
          AddMore(this.draftProps, val);
          this.draftProps.firstloading = false;
        });
      }
    }
  }

  public ReachedButtom() {
    if (this.tab === "tab-1") {
      ElmReachedButtom(this, ({ page, pageSize }) =>
        GetOwnArticles({ page: page, pageSize: pageSize, draft: false })
      );
    } else if (this.tab === "tab-2") {
      ElmReachedButtom(this.draftProps, ({ page, pageSize }) =>
        GetOwnArticles({ page: page, pageSize: pageSize, draft: true })
      );
    }
  }
}
</script>