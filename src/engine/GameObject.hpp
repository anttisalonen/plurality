#ifndef PLURALITY_ENGINE_GAMEOBJECT_HPP
#define PLURALITY_ENGINE_GAMEOBJECT_HPP

#include <string>
#include <vector>

#include "Component.hpp"

class GameObject {
	public:
		GameObject(const std::string& name) : mName(name) { }
		void addComponent(ComponentPtr c) { mComponents.push_back(c); }
		std::vector<ComponentPtr> getComponents() { return mComponents; }

	private:
		std::string mName;
		std::vector<ComponentPtr> mComponents;
};


#endif

