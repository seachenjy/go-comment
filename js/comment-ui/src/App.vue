<template>
  <div>
    <CommentForm></CommentForm>
    
  </div>
</template>

<script>
import CommentForm from './components/CommentForm'
import CommentItem from './components/CommentItem'

export default {
  name: 'App',
  components: {
    CommentForm,
    CommentItem
  }
}
</script>
