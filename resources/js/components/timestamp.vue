<template lang="pug">
  small.text-muted
    | {{ timestamp | moment('timezone', 'Asia/Jakarta').format('DD/MM/YYYY') }}
    translate at
    | &nbsp;
    | {{ timestamp | moment('timezone', 'Asia/Jakarta').format('HH:mm') }}
</template>

<script>
export default {
  name: 'Timestamp',
  props: ['timestamp'],
}
</script>
