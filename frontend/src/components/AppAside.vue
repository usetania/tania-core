<template>
  <aside id="aside">
    <div class="brand-image text-center">
      <img src="../assets/logo.png" alt="Tania logo" />
    </div>
    <b-nav vertical>

    </b-nav>
    <router-link to="/">Dashboard</router-link> |
    <router-link to="/about">About</router-link>
  </aside>
</template>

<style lang="scss" scoped>
aside#aside {
  background-color: #513969;
  height: 100%;

  .brand-image {
    padding: 15px 0;
    margin-bottom: 30px;

    img {
      width: 100px;
    }
  }
}
</style>
