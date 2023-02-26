#!/bin/bash

_episodes_completion() {
  actions="download list"
  COMPREPLY=($(compgen -W "$actions" "${COMP_WORDS[${COMP_CWORD}]}"))

  if [ ${COMP_CWORD} -gt 2 ]; then
    # check flags here
    action=${COMP_WORDS[2]}
    if [ "$action" = "download" ]; then
      COMPREPLY=($(compgen -W "<episode_list> -h" "${COMP_WORDS[${COMP_CWORD}]}"))
    else
      COMPREPLY=()
    fi
  fi
}

_music_completion() {
  actions="download list"
  COMPREPLY=($(compgen -W "$actions" "${COMP_WORDS[${COMP_CWORD}]}"))

  if [ ${COMP_CWORD} -gt 2 ]; then
    # check flags here
    action=${COMP_WORDS[2]}
    if [ "$action" = "download" ]; then
      COMPREPLY=($(compgen -W "<music_id_list> -h" "${COMP_WORDS[${COMP_CWORD}]}"))
    else
      COMPREPLY=()
    fi
  fi
}

_adcl_completions() {
  subcommands="episodes music"
  COMPREPLY=($(compgen -W "$subcommands" "${COMP_WORDS[1]}"))

  if [ ${COMP_CWORD} -gt 1 ]; then
    command=${COMP_WORDS[1]}

    case "$command" in
      "episodes")
          _episodes_completion
          ;;
      "music")
          _music_completion
          ;;
    esac


  fi
}

complete -F _adcl_completions adcl


