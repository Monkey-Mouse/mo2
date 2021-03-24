<template>
  <v-tooltip bottom content-class="">
    <template v-slot:activator="{ on, attrs }">
      <v-avatar v-bind="attrs" v-on="on" :size="size" color="brown">
        <v-img
          v-if="
            user.settings && user.settings.avatar && user.settings.avatar !== ''
          "
          :src="user.settings.avatar + '~thumb'"
        ></v-img>
        <span v-else class="white--text headline">{{ initials }}</span>
      </v-avatar>
    </template>
    <v-card-title>{{ user.name }}</v-card-title>
    <v-divider />
    <v-card-subtitle>{{
      email ? user.email : "github oauth account"
    }}</v-card-subtitle>
  </v-tooltip>
</template>

<script lang="ts">
import { User } from "@/models";
import { GetInitials, UserRole } from "../utils/index";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
@Component({
  components: {},
})
export default class Avatar extends Vue {
  @Prop()
  user!: User;
  @Prop()
  size?: number;
  get email() {
    return this.user.email.indexOf("@") > 0;
  }
  get isUser() {
    return this.user.roles && this.user.roles.indexOf(UserRole) >= 0;
  }

  get initials(): string {
    try {
      return GetInitials(this.user.name);
    } catch (error) {
      return "A";
    }
  }
}
</script>
