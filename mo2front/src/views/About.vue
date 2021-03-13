<template>
  <v-container class="fill-height">
    <v-row justify="center">
      <v-img
        alt="Logo"
        class="shrink mr-2"
        contain
        src="../assets/logo.png"
        transition="scale-transition"
        width="150"
      />
    </v-row>
    <v-row justify="center">
      <v-col class="text-center text-h2" cols="12">
        全心为用户打造的博客
      </v-col>
    </v-row>
    <v-row>
      <v-divider />
    </v-row>
    <v-row justify="center">
      <v-col class="" cols="12" offset-lg="2" lg="4">
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title class="text-h3"
              ><v-icon size="64">mdi-creation</v-icon>美丽</v-list-item-title
            >
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title class="text-h3">
              <v-icon size="64" @click="gotoGithub">mdi-github</v-icon
              >开源</v-list-item-title
            >
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title class="text-h3"
              ><v-icon size="64">mdi-bike-fast</v-icon>低延迟</v-list-item-title
            >
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title class="text-h3"
              ><v-icon size="64">mdi-alphabetical-variant-off</v-icon
              >无广告</v-list-item-title
            >
          </v-list-item-content>
        </v-list-item>
      </v-col>
      <v-divider v-if="!$vuetify.breakpoint.mobile" vertical />
      <v-col cols="12" lg="5" class="text">
        <h1>我们创造Mo2的原因</h1>
        <p>你们有没有思考过，一个理想中的博客网站的形态？</p>
        <p>
          首先没有创造的网站是没有价值的，
          一个博客网站的价值几乎完全来源于网站上发表博客的创作者。
          然而目前很多的博客网站却并没有给予创作者应有的特殊地位。
          创作者可能还是要付费才能免广告，有时要买会员才能下载别人的资源。
        </p>
        <p>
          同时，没有读者的创作是悲哀的。好文章的阅读，有时是思想的交互，有时是灵魂的共鸣。
          任何阅读过程都应该是神圣不可侵犯的，我们绝不会允许广告和注册提示之类的弹窗打断用户的阅读过程。
        </p>
        <p>
          因此我们的目标是：打造一个互联网上体验最好的博客平台，让作者爱上创作，让读者能投入阅读。<b
            >MAKE BLOGSITE GREAT AGAIN!</b
          >
        </p>
      </v-col>
    </v-row>
    <v-row>
      <v-divider />
    </v-row>
    <v-row justify="center">
      <v-col class="text-center" cols="12"
        ><h1>
          <a href="https://github.com/Monkey-Mouse/mo2"> CONTRIBUTING </a>
        </h1></v-col
      >
      <v-col lg="8" offset-lg="1" cols="12" class="text">
        <p>对，你们没有看错，contributing。</p>
        <p>
          创建一个超级好用的博客网站一直是我们团队的梦想，但是我们知道以我们目前的力量难以完成这个任务。
          并且，我们目前没有专门学习用户调研或者设计的成员。
        </p>
        <p>
          在我写这个About页面的时候，我们核心团队还只有两个人，一共进行了一个月不到的开发。我们在开发过程中想到：
          既然项目是开源的，为何不邀请我们有能力的用户帮助我们一起完善项目呢？我相信这也正是开源的意义所在。
        </p>
        <p>
          帮助我们并不仅仅意味着给我们提交pr，开启一个issue来汇报存在的bug或者提议项目以后的功能也大有助益，
          而一点点Sponsorship更是可能会让我们晚上开心到睡不着觉。
        </p>
        <p>我们欢迎并期待你们的任何支持！</p>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
@Component({
  components: {},
})
export default class About extends Vue {
  gotoGithub() {
    window.open("https://github.com/Monkey-Mouse/mo2", "_blank");
  }
}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
