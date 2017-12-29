<template lang="pug">
  aside#aside.app-aside.hidden-xs.bg-dark
    .aside-wrap
      .navi-wrap
        nav.navi.clearfix
          ul.nav
            li
              a.auto(href="")
                i.fa.fa-home.icon.text-primary
                span.font-bold.text-primary Dashboar
            li
              a.auto(href="")
                span.pull-right.text-muted
                  i.fa.fa-fw.fa-angle-right.text
                  i.fa.fa.fw.fa-angle-down.text-active
                i.fa.fa-leaf
                span Production

</template>

<script>
export default {
  name: 'AppAsideComponent'
}
</script>
