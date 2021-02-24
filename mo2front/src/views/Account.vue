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
      <blog-time-line-list :blogs="blogs" />
    </v-container>
  </div>
</template>

<script lang="ts">
import { BlogBrief, User } from "@/models";
import { Copy, GetOwnArticles, GetUserArticles, GetUserData } from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import Avatar from "../components/UserAvatar.vue";
@Component({
  components: {
    BlogTimeLineList,
    Avatar,
  },
})
export default class Account extends Vue {
  @Prop()
  user!: User;
  displayUser: User;
  uid!: string;
  blogs: BlogBrief[] = [];
  created() {
    this.uid = this.$route.params["id"];
    if (this.uid === undefined || this.uid === this.user.id) {
      this.uid = this.user.id;
      this.displayUser = Copy(this.user);
      GetOwnArticles({ page: 0, pageSize: 10, draft: false }).then((data) => {
        this.blogs = data;
      });
    } else {
      GetUserData(this.uid).then((u) => {
        this.displayUser = u;
        GetUserArticles({
          page: 0,
          pageSize: 10,
          draft: false,
          id: this.uid,
        }).then((data) => {
          this.blogs = data;
        });
      });
    }
  }
}
</script>