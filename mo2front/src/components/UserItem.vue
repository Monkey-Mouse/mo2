<template>
  <v-card :to="'/account/' + user.id" class="mx-auto" width="300px" outlined>
    <v-list-item three-line>
      <user-avatar :user="user" :enableEdit="false"></user-avatar>
      <v-list-item-content class="ml-1">
        <v-list-item-title class="text-h5 mb-1">
          {{ user.name }}
        </v-list-item-title>
        <v-list-item-subtitle>{{ user.email }}</v-list-item-subtitle>
      </v-list-item-content>
    </v-list-item>
  </v-card>
</template>

<script lang="ts">
import { User, UserListData } from "@/models";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import UserAvatar from "./UserAvatar.vue";
@Component({
  components: {
    UserAvatar,
  },
})
export default class UserItem extends Vue {
  @Prop()
  user!: UserListData;
}
</script>

<style>
</style>