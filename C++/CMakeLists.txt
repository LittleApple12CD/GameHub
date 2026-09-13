cmake_minimum_required(VERSION 3.28)
project(GameHub LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# 可执行文件统一输出到项目根目录下的 dist/
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/dist)

find_package(SFML 3 REQUIRED COMPONENTS Graphics Window System)

set(SOURCES
    src/main.cpp
    src/Menu.cpp
    src/games/SnakeGame.cpp
    src/games/TetrisGame.cpp
    src/games/MinesweeperGame.cpp
    src/games/FlappyGame.cpp
    src/games/TankGame.cpp
    src/games/BreakoutGame.cpp
    src/games/TicTacToeGame.cpp
)

add_executable(GameHub ${SOURCES})

target_link_libraries(GameHub PRIVATE
    SFML::Graphics
    SFML::Window
    SFML::System
)

if(WIN32 AND MINGW)
    set(RESOURCE_FILE ${CMAKE_SOURCE_DIR}/resources.rc)
    set(RESOURCE_OBJ  ${CMAKE_BINARY_DIR}/resources.o)

    if(EXISTS ${RESOURCE_FILE})
        add_custom_command(
            OUTPUT ${RESOURCE_OBJ}
            COMMAND windres ${RESOURCE_FILE} -o ${RESOURCE_OBJ}
            DEPENDS ${RESOURCE_FILE}
            WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
            COMMENT "Compiling resources.rc with windres"
        )
        target_sources(GameHub PRIVATE ${RESOURCE_OBJ})
    else()
        message(WARNING "resources.rc not found, exe icon will not be set")
    endif()
endif()

if(EXISTS ${CMAKE_SOURCE_DIR}/assets)
    add_custom_command(TARGET GameHub POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E copy_directory
            ${CMAKE_SOURCE_DIR}/assets
            ${CMAKE_RUNTIME_OUTPUT_DIRECTORY}/assets
        COMMENT "Copying assets to dist/"
    )
endif()