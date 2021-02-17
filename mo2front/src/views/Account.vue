<template>
  <div>
    <v-parallax
      dark
      src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"
      height="400"
    >
      <v-row align="center" justify="center">
        <v-col class="text-center" cols="12">
          <v-img class="v-avatar" :src="displayUser.avatar"></v-img>
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
import { Copy, GetUserData } from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
@Component({
  components: {
    BlogTimeLineList,
  },
})
export default class Account extends Vue {
  @Prop()
  user!: User;
  displayUser: User;
  uid!: string;
  blogs: BlogBrief[] = Array<BlogBrief>(10).fill({
    id: "string",
    title: "MO2",
    cover: "https://picsum.photos/500/300?image=40",
    rate: 4.3,
    description:
      "Lorem ipsum dolor sit amet, no nam oblique veritus. Commune scaevola imperdiet nec ut, sed euismod convenire principes at. Est et nobis iisque percipit, an vim zril disputando voluptatibus, vix an salutandi sententiae.",
    createTime: "2021/2/9",
    author: "Leezeeyee",
  });
  created() {
    this.uid = this.$route.params["id"];
    if (this.uid === undefined) {
      this.uid = this.user.id;
      this.displayUser = Copy(this.user);
    } else {
      GetUserData(this.uid).then((u) => {
        this.displayUser = u;
      });
    }
  }
}
</script>