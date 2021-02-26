<template>
  <div>
    <v-parallax
      dark
      src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"
      height="400"
    >
      <v-row align="center" justify="center">
        <v-col class="text-center" cols="12">
          <avatar :size="80" :user="displayUser" />
          <!-- <v-img class="v-avatar" :src="displayUser.avatar"></v-img> -->
          <h1 class="display-1 font-weight-thin mb-4">
            {{ displayUser.name }}
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

        <v-tab v-if="ownPage" @click="loadDraft" href="#tab-2">
          DraftBox
        </v-tab>
      </v-tabs>

      <v-tabs-items v-model="tab">
        <v-tab-item :value="'tab-1'">
          <v-card flat>
            <blog-time-line-list v-if="!firstloading" :blogs="blogs" />
            <blog-skeleton v-if="loading" :num="pagesize" />
          </v-card>
        </v-tab-item>
        <v-tab-item :value="'tab-2'">
          <v-card flat>
            <blog-time-line-list
              v-if="!draftProps.firstloading"
              :blogs="draftProps.blogs"
              :draft="true"
            />
            <blog-skeleton v-if="draftProps.loading" :num="pagesize" />
          </v-card>
        </v-tab-item>
      </v-tabs-items>
    </v-container>
  </div>
</template>

<script lang="ts">
import { BlogBrief, User } from "@/models";
import {
  AddMore,
  BlogAutoLoader,
  Copy,
  ElmReachedButtom,
  GetOwnArticles,
  GetUserArticles,
  GetUserData,
} from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import Avatar from "../components/UserAvatar.vue";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
@Component({
  components: {
    BlogTimeLineList,
    Avatar,
    BlogSkeleton,
  },
})
export default class Account extends Vue implements BlogAutoLoader {
  @Prop()
  user!: User;
  displayUser: User;
  uid!: string;
  blogs: BlogBrief[] = [];
  loading = true;
  firstloading = true;
  page = 0;
  pagesize = 5;
  nomore = false;
  tab = 1;
  ownPage = false;

  draftProps: BlogAutoLoader = {
    loading: true,
    firstloading: true,
    page: 0,
    pagesize: 5,
    nomore: false,
    blogs: [],
    ReachedButtom: () => {},
  };
  created() {
    this.uid = this.$route.params["id"];
    if (this.uid === undefined || this.uid === this.user.id) {
      this.uid = this.user.id;
      this.displayUser = this.user;
      this.ownPage = true;
      GetOwnArticles({
        page: this.page++,
        pageSize: this.pagesize,
        draft: false,
      }).then((data) => {
        AddMore(this, data);
        this.firstloading = false;
      });
    } else {
      GetUserData(this.uid).then((u) => {
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
      });
    }
  }
  @Watch("user")
  userChange() {
    if (
      this.uid === undefined ||
      this.uid === "" ||
      this.uid === this.user.id
    ) {
      this.uid = this.user.id;
      this.displayUser = this.user;
    }
  }
  loadDraft() {
    if (this.draftProps.firstloading) {
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
    if (this.tab === 1) {
      ElmReachedButtom(this, ({ page, pageSize }) =>
        GetOwnArticles({ page: page, pageSize: pageSize, draft: false })
      );
    } else {
      ElmReachedButtom(this.draftProps, ({ page, pageSize }) =>
        GetOwnArticles({ page: page, pageSize: pageSize, draft: true })
      );
    }
  }
}
</script>